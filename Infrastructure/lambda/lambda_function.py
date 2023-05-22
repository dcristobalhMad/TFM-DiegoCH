import re
import base64
import json


def parse_log(log):
    # Define the regular expression pattern
    pattern = r'^(\S+) (\S+) \[(.*?)\] (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) (\S+) \{(.*?)\} "(.*?)"$'

    # Match the pattern against the log string
    match = re.match(pattern, log)

    if match:
        # Extract the desired information from the matched groups
        ip1 = match.group(1)  # 127.0.0.1:4398
        ip2 = match.group(2)  # 127.0.0.1:80
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
        user_agent = match.group(14)  # localhost|curl/7.29.0
        request = match.group(15)  # GET / HTTP/1.1

        # Create a dictionary with the extracted fields
        log_fields = {
            "IP1": ip1,
            "IP2": ip2,
            "Timestamp": timestamp,
            "Section": section,
            "Resource": resource,
            "Values": values,
            "Status Code": status_code,
            "Size": size,
            "Dash1": dash1,
            "Dash2": dash2,
            "Flags": flags,
            "Values2": values2,
            "Values3": values3,
            "User Agent": user_agent,
            "Request": request,
        }

        return log_fields
    else:
        return None


def lambda_handler(event, context):
    output_records = []
    for record in event["records"]:
        payload = base64.b64decode(record["data"]).decode("utf-8")
        # Parse each line of the payload
        parsed_data = []
        # Load string as a json object
        payload_json = json.loads(payload)
        # Get the message from the json object
        payload = payload_json["message"]
        for line in payload.splitlines():
            parsed_log = parse_log(line)
            print("Parsed line:")
            print(parsed_log)
            if parsed_log:
                parsed_data.append(parsed_log)
                print("Parsed data:")
                print(parsed_data)

        # Convert parsed data to JSON and add to output records
        output_payload = json.dumps(parsed_data)
        output_records.append(
            {
                "recordId": record["recordId"],
                "result": "Ok",
                "data": base64.b64encode(output_payload.encode("utf-8")).decode(
                    "utf-8"
                ),
            }
        )
        print("Output payload:")
        print(output_payload)
        print("Output records:")
        print(output_records)
    return {"records": output_records}
