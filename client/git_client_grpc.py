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
import json

import grpc
import gitapi_pb2
import gitapi_pb2_grpc

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
            callback=git_process)
        
        syntax_commit = (Syntax(
            'commit', 
            help='git commit commits the startup config in github')
            .add_unnamed_argument('commitMessage',
                help='commitMessage that is added to the commit'))
        commit = git.add_command(syntax_commit, 
            update_location=False,
            callback=git_process)

        syntax_pull_request = (Syntax(
            'pull-request', 
            help='git pull-request creates a pull request based on the commits in github')
            .add_unnamed_argument('prMessage',
                help='pull request message that is added to the pull-request'))
        pull_request = git.add_command(syntax_pull_request, 
            update_location=False,
            callback=git_process)

def git_process(state, output, arguments, **_kwargs):
    #print(state.server_data_store)
    print(arguments)

    if arguments == 'branch':
        print(arguments)
    elif arguments == 'commit':
        print(arguments)
    elif arguments == 'pull-request':
        print(arguments)
    else:
        print(arguments)
        #output.print_error_line('invalid argument; use branch, commit, pull-request')

    # path = build_path('/system/aaa/authentication/session')
    # data = state.server.get_data_store(
    #     DataStore.State).get_data(
    #         path,
    #         recursive=True,
    #         add_missing_containers=True)
    # #print(data)
    # for session in data.get_descendants('/system/aaa/authentication/session'):
    #     print(session)

    # path = build_path('/git-client')
    # data = state.server_data_store.get_data(
    #         path, 
    #         recursive=True,
    #         include_container_children=True)

    #print(server_data.get_descendants('/git-client'))

    # for git in data.get_descendants('/git-client'):
    #     print(git)

    path = build_path('/git-client')
    result = state.server_data_store.get_json(
        path, 
        recursive=True,
        include_field_defaults=True)
    #print(result)

    result_json_obj = json.loads(result)

    print(result_json_obj)

    parent_path = str(state.line_commands.path)
    print(parent_path)

    with grpc.insecure_channel('localhost:7777') as channel:
        client = gitapi_pb2_grpc.GitStub(channel)
        a = gitapi_pb2.Action(
            kind="branch",
            attributes="Wim commit"
            )

        response = client.Command(a)
        print(response.response)

    #output.print_error_line('Not sure if all paramters exist')