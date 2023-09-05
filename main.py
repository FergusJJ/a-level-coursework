from src.default_messages import Messages
from src.menus.main_menu import Menu, return_site_name
from src.colors import Style
from src.img import Img
from src.utils.init_tasks import get_taskdata, pyToGo
from colorama import init
from time import sleep
import os
from time import strftime
# https://www.delftstack.com/howto/python/python-clear-console/
clearConsole = lambda:print(chr(27) + "[2J")

def start_prog():
    
    init()
    
    start_sites()
    
def start_sites():
    selected_site = "-1"
    while selected_site == "-1":
        clearConsole()

        Img.show_logo()
        Menu.show_ver()
        Menu.show_sites()
        selected_site = Menu.get_choice()

    task_list = get_taskdata(selected_site)
    if len(task_list) == 0:
        Messages.no_profiles_for_site()
        sleep(1)
        start_sites()
    pyToGo(task_list)
    await_user_input()
    clearConsole()
    
def await_user_input():
    timestamp = strftime('%H:%M:%S')
    _ = input(f"[{timestamp}] Press enter to continue")


if __name__ == "__main__":
    try:
        while True:
            start_prog() 
    except KeyboardInterrupt:
        return_site_name(0)
