from src.colors import Style
from art import tprint

import sys

class Img:

    @staticmethod
    def show_logo():
        sys.stdout.write(Style.MAGENTA)
        tprint("DeHype",font="smslant")
        sys.stdout.write(Style.RESET)