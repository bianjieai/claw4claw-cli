import os, subprocess, re

def run_cmd(cmd):
    try:
        return subprocess.check_output(cmd, shell=True, text=True).strip()
    except Exception as e:
        print(f"Error executing command: {cmd}\n{e}")
        return ''

ref_name = os.environ.get('GITHUB_REF_NAME', '')
repo = os.environ.get('GITHUB_REPOSITORY', '')

if not ref_name or not repo:
    print("GITHUB_REF_NAME or GITHUB_REPOSITORY environment variables are not set.")
    exit(1)

prev_tag = run_cmd(f'git describe --tags --abbrev=0 {ref_name}^')
if prev_tag:
    log_range = f'{prev_tag}..{ref_name}'
    compare_url = f'https://github.com/{repo}/compare/{prev_tag}...{ref_name}'
    baseline_str = f'（对比 {prev_tag}）'
else:
    log_range = ref_name
    compare_url = f'https://github.com/{repo}/commits/{ref_name}'
    baseline_str = ''

commits = run_cmd(f'git log {log_range} --oneline').split('\n')

features, fixes, chores = [], [], []

type_map = {
    'feat': '新特性',
    'perf': '性能优化',
    'fix': '修复',
    'bug': '修复',
    'docs': '文档',
    'chore': '杂项',
    'refactor': '重构',
    'ci': 'CI流程',
    'test': '测试',
    'build': '构建'
}

for commit in commits:
    if not commit: continue
    parts = commit.split(' ', 1)
    hash_short = parts[0]
    msg = parts[1] if len(parts) > 1 else ''
    commit_url = f'https://github.com/{repo}/commit/{hash_short}'
    
    match = re.match(r'^([a-zA-Z]+)(?:\(([^)]+)\))?:\s*(.*)', msg)
    if match:
        c_type, scope, desc = match.group(1), match.group(2), match.group(3)
        if not scope:
            scope = type_map.get(c_type, '综合')
        item = f'- **{scope}**：{desc} ([{hash_short}]({commit_url}))'
        if c_type in ['feat', 'perf']: features.append(item)
        elif c_type in ['fix', 'bug']: fixes.append(item)
        else: chores.append(item)
    else:
        chores.append(f'- **其他**：{msg} ([{hash_short}]({commit_url}))')

md = [f'# 🎉 Release {ref_name}\n\n本次发布{baseline_str}主要包含以下特性的更新和问题的修复：\n']
if features: md.extend(['## ✨ 新特性与优化 (Features & Enhancements)'] + features + [''])
if fixes: md.extend(['## 🐛 问题修复 (Bug Fixes)'] + fixes + [''])
if chores: md.extend(['## 📝 文档与杂项 (Documentation & Chores)'] + chores + [''])

md.extend(['---', f'**🔍 完整代码变更**: [{prev_tag}...{ref_name}]({compare_url})' if prev_tag else f'**🔍 完整代码变更**: [查看变更]({compare_url})'])

with open('release_note.md', 'w', encoding='utf-8') as f:
    f.write('\n'.join(md))

print("Release note generated successfully at release_note.md")
