#!/usr/bin/python
###########################################################################
# Description:
#
# Copyright (c) 2020 Nokia
###########################################################################
from srlinux.location.build_path import build_path
from srlinux.mgmt.cli.tools_plugin import ToolsPlugin
from srlinux.mgmt.cli.required_plugin import RequiredPlugin
from srlinux.schema.data_store import DataStore
from srlinux.syntax import Syntax
from srlinux.mgmt.cli.parse_error import ParseError

from srlinux.schema.data_store import DataStore
from srlinux.location import build_path
import requests
import json

class Plugin(ToolsPlugin):
    '''
        git
    '''

    def get_required_plugins(self):
        return [
            RequiredPlugin('tools_mode')
        ]
    def on_tools_load(self, state):
        syntax_git = Syntax('git', 
            help='`Git interaction as a client`')
        git = state.command_tree.tools_mode.root.add_command(syntax_git,
            update_location=False)

        syntax_branch = Syntax('branch', 
            help='git branch creates a branch in github')
        branch = git.add_command(syntax_branch, 
            update_location=False, 
            callback=git_branch_process)
        
        syntax_commit = (Syntax(
            'commit', 
            help='git commit commits the startup config in github')
            .add_unnamed_argument('commitMessage',
                help='commitMessage that is added to the commit'))
        commit = git.add_command(syntax_commit, 
            update_location=False,
            callback=git_commit_process)

        syntax_pull_request = (Syntax(
            'pull-request', 
            help='git pull-request creates a pull request based on the commits in github')
            .add_unnamed_argument('prMessage',
                help='pull request message that is added to the pull-request'))
        pull_request = git.add_command(syntax_pull_request, 
            update_location=False,
            callback=git_pullrequest_process)

def git_branch_process(state, output, arguments, **_kwargs):
    #print(state.server_data_store)
    #print(arguments)

    path = build_path('/git-client')
    result = state.server_data_store.get_json(
        path, 
        recursive=True,
        include_field_defaults=True)

    result_json_obj = json.loads(result)
    #print(result_json_obj)

    #parent_path = str(state.line_commands.path)
    #print(parent_path)

    url = "http://localhost:7777"
    payload = {
        "method": "Server.Branch",
        "params": [{"Comment": "H"}],
        "jsonrpc": "2.0",
        "id": 0,
    }
    response = requests.post(url, json=payload).json()

    if response["result"] != 'success':
        output.print_error_line(response["result"])
    assert response["id"] == 0

def git_commit_process(state, output, arguments, **_kwargs):
    #print(state.server_data_store)
    #print(arguments)
    comment = arguments.get('commitMessage')

    path = build_path('/git-client')
    result = state.server_data_store.get_json(
        path, 
        recursive=True,
        include_field_defaults=True)

    result_json_obj = json.loads(result)
    #print(result_json_obj)

    #parent_path = str(state.line_commands.path)
    #print(parent_path)

    url = "http://localhost:7777"
    payload = {
        "method": "Server.Commit",
        "params": [{"Comment": comment}],
        "jsonrpc": "2.0",
        "id": 0,
    }
    response = requests.post(url, json=payload).json()

    if response["result"] != 'success':
        output.print_error_line(response["result"])
    assert response["id"] == 0

def git_pullrequest_process(state, output, arguments, **_kwargs):
    #print(state.server_data_store)
    #print(arguments)
    comment = arguments.get('prMessage')

    path = build_path('/git-client')
    result = state.server_data_store.get_json(
        path, 
        recursive=True,
        include_field_defaults=True)

    result_json_obj = json.loads(result)
    #print(result_json_obj)

    #parent_path = str(state.line_commands.path)
    #print(parent_path)

    url = "http://localhost:7777"
    payload = {
        "method": "Server.PullRequest",
        "params": [{"Comment": comment}],
        "jsonrpc": "2.0",
        "id": 0,
    }
    response = requests.post(url, json=payload).json()

    if response["result"] != 'success':
        output.print_error_line(response["result"])
    assert response["id"] == 0