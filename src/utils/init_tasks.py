import csv
from enum import Flag
import os 
import ctypes
import re

from src.default_messages import Messages

class GoSlice(ctypes.Structure):
    _fields_ = [("data", ctypes.POINTER(ctypes.c_void_p)),
                ("len", ctypes.c_longlong), ("cap", ctypes.c_longlong)]


so = ctypes.CDLL("gosrc/pyconverter.so")
so.convertToGo.argtypes = [GoSlice]


def get_taskdata(site: str) -> list:
    folder = os.listdir("profiles/")
    taskList = []
    for files in folder:
        file_path = f"profiles/{files}"
        with open(file_path, newline="") as file:
            csv_file = csv.DictReader(file,delimiter=',', quotechar='"', quoting=csv.QUOTE_MINIMAL)
            for row in csv_file:
                store = row["store"]
                if store.lower() == site:
                    taskList.append(row)
            file.close()

    return taskList

def pyToGo(tasklist: list):
    regex = r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b' #https://www.geeksforgeeks.org/check-if-email-address-valid-or-not-in-python/
    flag = False
    for task in tasklist:
        args = []
        #check for invalid values
        if (task["mode"].lower()).strip() != "desktop" or (task["mode"].lower()).strip() != "mobile":
            task["mode"] = "desktop"
        else:
            task["mode"] = (task["mode"].lower()).strip()

        if "https" not in (task["product"].lower()).strip():
            Messages.invalid_url(task["email"],task["product"])
            #wont be a valid url so not going to add to the task list
            continue
        try:
            float(task["size"]) # will hit except if its not a number
        except ValueError:
            Messages.invalid_size(task["email"],task["size"])
            continue

        #size 0 == random
        task["size"] = task["size"].strip()
        if "." not in task["size"]:
            temp_size = int(task["size"])
            if temp_size == 0:
                task["size"] = "random" 
        
        if len(task["proxy"]) == 0:
            task["proxy"] = "localhost"

        if len(task["delay"]) == 0:
            Messages.invalid_delay(task["email"])
            print(chr(27) + "[2J")
            task["delay"] = 1
        
        try:
            float(task["delay"]) # will hit except if its not a number
        except ValueError:
            Messages.invalid_delay(task["email"],task["delay"])
            task["delay"] = 0


        if len(task["first name"]) == 0:
            Messages.invalid_name(task["email"])
            print(chr(27) + "[2J")
            continue
        if len(task["last name"]) == 0:
            Messages.invalid_name(task["email"])
            continue
        #email is invalid
        if not re.fullmatch(regex, task["email"]):
            Messages.invalid_email(task["email"])
            continue
        if len(task["line1"].strip()) == 0:
            Messages.invalid_line1(task["email"])
            continue
        if len(task["city"].strip()) == 0:
            Messages.invalid_city(task["email"])
            continue
        if len(task["postcode"].strip()) == 0:
            Messages.invalid_postcode(task["email"])
            continue
        if len(task["card number"].strip()) == 0:
            Messages.invalid_cc_num(task["email"])
            continue
        if len(task["expiry month"].strip()) < 2:
            Messages.invalid_cc_exp_month_year(task["email"])
            continue
        if len(task["expiry year"].strip()) < 2:
            Messages.invalid_cc_exp_month_year(task["email"])
            continue
        if len(task["cvc"].strip()) == 0:
            Messages.invalid_cvc(task["email"])
            continue
        flag = True
        for key in task:
            temp_string = str(task[key])
            args.append(ctypes.cast(ctypes.c_char_p(temp_string.encode("utf-8")),ctypes.c_void_p))

        #unpacks args into a c array
        temp_array = (ctypes.c_void_p * len(args))(*args)
        temp_go_slice = GoSlice(temp_array,len(temp_array),len(temp_array))
        so.convertToGo(temp_go_slice)
        
    if flag:

        so.checkProfileMap()
    






