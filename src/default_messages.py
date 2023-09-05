from src.colors import Style

import os
import sys
import time

class Messages:

    #timestamp = time.strftime('%H:%M:%S')

    @staticmethod
    def no_profiles():
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] COULDN'T FIND ANY .CSV FILES IN THE PROFILES FOLDER, PLEASE CREATE A FILE BEFORE USING THE BOT")
        time.sleep(1)
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def no_profiles_for_site():
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] COULDN'T FIND ANY PROFILES FOR THAT SITE, PLEASE CREATE PROFILES BEFORE USING THE BOT")
        time.sleep(1)
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def unknown_err():
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] AN UNKNOWN  ERROR HAS OCCURRED")
        time.sleep(1)
        sys.stdout.write(Style.RESET)
        os.system('cls')
    
    @staticmethod
    def bad_input():
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] THAT IS NOT AN OPTION, PLEASE TRY AGAIN")
        sys.stdout.write(Style.RESET)
        

    @staticmethod
    def invalid_url(profile_email, profile_url: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - INVALID URL \"{profile_url}\"")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_size(profile_email, profile_size: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - INVALID SIZE \"{profile_size}\"")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_delay(profile_email, profile_delay: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - INVALID DELAY, MUST NOT CONTAIN ANY LETTERS & MUST NOT BE LEFT BLANK \"{profile_delay}\"")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_name(profile_email):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - INVALID FIRSTNAME OR LASTNAME, CANNOT BE BLANK")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_email(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - INVALID EMAIL")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_line1(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - ADDRESS LINE 1 CANNOT BE BLANK")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_city(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - CITY CANNOT BE BLANK")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_city(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - POSTCODE CANNOT BE BLANK")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_cc_num(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - CREDIT CARD NUMBER CANNOT BE BLANK")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_cc_exp_month_year(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - INVALID EXPIRY DATE")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def invalid_cvc(profile_email: str):
        sys.stdout.write(Style.RED)
        timestamp = time.strftime('%H:%M:%S')
        print(f"[{timestamp}] SKIPPING {profile_email} - CVC")
        
        sys.stdout.write(Style.RESET)
        os.system('cls')

    @staticmethod
    def random_closing_msg(msg_num: int):
        if msg_num == 1:
            sys.stdout.write(Style.BLUE)
            timestamp = time.strftime('%H:%M:%S')
            print(f"[{timestamp}] Goodbye...")
            time.sleep(1)
            sys.stdout.write(Style.RESET)
            sys.exit()
        elif msg_num == 2:
            sys.stdout.write(Style.BLUE)
            timestamp = time.strftime('%H:%M:%S')
            print(f"[{timestamp}] Shutting down...")
            time.sleep(1)
            sys.stdout.write(Style.RESET)
            sys.exit()
        elif msg_num == 3:
            sys.stdout.write(Style.BLUE)
            timestamp = time.strftime('%H:%M:%S')
            print(f"[{timestamp}] Cutting the lights...")
            time.sleep(1)
            sys.stdout.write(Style.RESET)
            sys.exit()