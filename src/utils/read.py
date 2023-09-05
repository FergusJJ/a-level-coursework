import csv
class File_reader:


    @staticmethod
    def read_sites_from_csv(fname: str) -> int:
        #
        ZALANDO_COUNTER = 0 
        ASOS_COUNTER = 0
        FOOTLOCKER_COUNTER = 0
        BSTN_COUNTER = 0
        #
        file_path = f"profiles/{fname}"
        with open(file_path, newline="") as file:
            csv_file = csv.DictReader(file,delimiter=',', quotechar='"', quoting=csv.QUOTE_MINIMAL)

            for row in csv_file:
                store = row["store"]
                if store.lower() == "zalando":
                    ZALANDO_COUNTER += 1 
                elif store.lower() == "asos":
                    ASOS_COUNTER += 1
                elif store.lower() == "footlocker":
                    FOOTLOCKER_COUNTER += 1
                elif store.lower() == "bstn":
                    BSTN_COUNTER += 1

            file.close()
        return ZALANDO_COUNTER,ASOS_COUNTER,FOOTLOCKER_COUNTER,BSTN_COUNTER
             
            
            
