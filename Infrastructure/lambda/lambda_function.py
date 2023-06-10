import re
import base64
import json
import os


def parse_log(log):
    # Define the regular expression pattern
    pattern = r'^(\S+) (\S+) \[(.*?)\] (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) \{(.*?)\} "(.*?)"$'

    # Match the pattern against the log string
    match = re.match(pattern, log)
    if match:
        # Extract the desired information from the matched groups
        ip1 = match.group(1)  # 127.0.0.1:4398
        ip2 = match.group(2)  # 127.0.0.1:80 THIS IS THE SERVER IP
        timestamp = match.group(3)  # 22/May/2023:03:06:13.724
        section = match.group(4)  # main
        resource = match.group(5)  # wsx/ws5
        values = match.group(6)  # 0/0/0/4/4
        status_code = match.group(7)  # 200
        size = match.group(8)  # 237
        dash1 = match.group(9)  # -
        dash2 = match.group(10)  # -
        flags = match.group(11)  # --NN
        values2 = match.group(12)  # 1/1/0/0/0
        values3 = match.group(13)  # 0/0
        user_agent = match.group(14)
        request = match.group(15)  # GET / HTTP/1.1

        ###### ENCODE VALUES ######
        # Encode the values
        ip2 = base64.b64encode(ip2.encode("utf-8")).decode("utf-8")
        user_agent = base64.b64encode(user_agent.encode("utf-8")).decode("utf-8")

        ###### TRANSFORM VALUES ######
        # Transform the values
        size = int(size.strip())
        status_code = int(status_code.strip())
        timestamp = timestamp.strip()
        # Create a dictionary with the extracted fields
        log_fields = {
            "Client_IP": ip1,
            "Server_IP": ip2,
            "Timestamp": timestamp,
            "Virtual_Host": section,
            "Server": resource,
            "Stats": values,
            "Status_Code": status_code,
            "Response_Size": size,
            "Referrer": dash1,
            "Header_user_agent": dash2,
            "SSL_information": flags,
            "SSL_stats": values2,
            "Server_stats": values3,
            "User_Agent": user_agent,
            "HTTP_Request": request,
        }

        return log_fields
    else:
        return None


def lambda_handler(event, context):
    output_records = []
    for record in event["records"]:
        payload = base64.b64decode(record["data"]).decode("utf-8")
        # Load string as a json object
        payload_json = json.loads(payload)
        # Get the message from the json object
        payload = payload_json["message"]
        parsed_log = parse_log(payload)
        output_payload = json.dumps(parsed_log)
        output_records.append(
            {
                "recordId": record["recordId"],
                "result": "Ok",
                "data": base64.b64encode(output_payload.encode("utf-8")).decode(
                    "utf-8"
                ),
            }
        )
    return {"records": output_records}
