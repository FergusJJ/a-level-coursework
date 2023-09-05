from src.colors import Style
from src.default_messages import Messages
from src.utils.read import File_reader

##
from random import randrange
import sys
import os
import time
##
class Menu:

    @staticmethod
    def get_profiles():
        ZALANDO_COUNTER = 0 
        ASOS_COUNTER = 0
        FOOTLOCKER_COUNTER = 0
        BSTN_COUNTER = 0
        #gets files in site folder
        folder = os.listdir("profiles/")
        if len(folder) == 0:
            Messages.no_profiles()
            import main
            main.start_prog()
        elif len(folder) == 1:
            selected_file = folder[0]
            returned_zalando_counter, returned_asos_counter, returned_footlocker_counter, returned_bstn_counter = File_reader.read_sites_from_csv(selected_file)
            ZALANDO_COUNTER += returned_zalando_counter
            ASOS_COUNTER += returned_asos_counter
            FOOTLOCKER_COUNTER += returned_footlocker_counter
            BSTN_COUNTER += returned_bstn_counter
        elif len(folder) > 1:
            for i in folder:
                selected_file = i
                returned_zalando_counter, returned_asos_counter, returned_footlocker_counter, returned_bstn_counter = File_reader.read_sites_from_csv(selected_file)
                ZALANDO_COUNTER += returned_zalando_counter
                ASOS_COUNTER += returned_asos_counter
                FOOTLOCKER_COUNTER += returned_footlocker_counter
                BSTN_COUNTER += returned_bstn_counter
            #do stuff but for multiple
        else:
            Messages.unknown_err()
            import main
            main.start_prog()

        return ZALANDO_COUNTER, ASOS_COUNTER, FOOTLOCKER_COUNTER, BSTN_COUNTER

    @staticmethod
    def show_ver():
        sys.stdout.write(Style.YELLOW)
        print("\nWelcome || Version 0.0.1")
        sys.stdout.write(Style.RESET)
    
    @staticmethod
    def show_sites():
        ZALANDO_COUNTER, ASOS_COUNTER, FOOTLOCKER_COUNTER, BSTN_COUNTER = Menu.get_profiles()
        sys.stdout.write(Style.YELLOW)
        print("\nSitelist:")
        sys.stdout.write(Style.RESET)
        print(f"[ 1 || Zalando || {ZALANDO_COUNTER} Profiles ]")
        print(f"[ 2 || Asos || {ASOS_COUNTER} Profiles ]")
        print(f"[ 3 || BSTN || {BSTN_COUNTER} Profiles ]")
        print(f"[ 4 || Footlocker || {FOOTLOCKER_COUNTER} Profiles ]")
        
        sys.stdout.write(Style.RED)
        print ("\n[ 0 || Exit ]\n")
        sys.stdout.write(Style.RESET)

    @staticmethod
    def get_choice():
        try:
            sys.stdout.write(Style.YELLOW)
            print("Which site would you like to bot?")
            sys.stdout.write(Style.RESET)
            choice = input("> ")
            site_choice = return_site_name(int(choice))
            return site_choice
            

        except KeyboardInterrupt:
            msg_num = randrange(1,4)
            Messages.random_closing_msg(msg_num)

        except ValueError:
            Messages.bad_input()
            return "-1"
            



def return_site_name(choice: int) -> str:
    if choice == 0:
        sys.stdout.write(Style.BLUE)
        print("GOODYBYE...")
        sys.stdout.write(Style.RESET)
        time.sleep(1)
        sys.exit()

    elif choice == 1:
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] Starting zalando")
        return "zalando"
            

    elif choice == 2:
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] Starting Asos")
        return "asos"
            

    elif choice == 3:
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] Starting BSTN")
        return "bstn"
            

    elif choice == 4:
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] Starting Footlocker")
        return "footlocker"
    else:
        return "-1"

    
        