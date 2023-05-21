import re
import base64
import json


def parse_line(line):
    pattern = r'^(\w+\s\d+\s\d{2}:\d{2}:\d{2})\s\w+\[(\d+)\]:\s([\d.]+:\d+)\s([\d.]+:\d+)\s\[(\d+\/\w+\/\d+:\d{2}:\d{2}:\d{2}\.\d+)\]\s(\w+)\s(\w+)\/(\w+)\s(\d+\/\d+\/\d+\/\d+\/\d+)\s(\d+)\s(\d+)\s-\s-\s--\w+\s(\d+\/\d+\/\d+\/\d+\/\d+)\s(\d+\/\d+)\s\{(.*?)\}\s"(.*?)"$'
    match = re.match(pattern, line)
    if match:
        timestamp = match.group(1)
        process_id = match.group(2)
        source_address = match.group(3)
        destination_address = match.group(4)
        request_timestamp = match.group(5)
        frontend_name = match.group(6)
        backend_name = match.group(7)
        server_name = match.group(8)
        timings = match.group(9)
        status_code = match.group(10)
        bytes_read = match.group(11)
        connection_times = match.group(12)
        session_times = match.group(13)
        user_agent = match.group(14)
        request = match.group(15)

        return {
            "timestamp": timestamp,
            "process_id": process_id,
            "source_address": source_address,
            "destination_address": destination_address,
            "request_timestamp": request_timestamp,
            "frontend_name": frontend_name,
            "backend_name": backend_name,
            "server_name": server_name,
            "timings": timings,
            "status_code": status_code,
            "bytes_read": bytes_read,
            "connection_times": connection_times,
            "session_times": session_times,
            "user_agent": user_agent,
            "request": request,
        }
    else:
        return None


def lambda_handler(event, context):
    output_records = []

    for record in event["records"]:
        payload = base64.b64decode(record["data"]).decode("utf-8")

        # Parse each line of the payload
        parsed_data = []
        for line in payload.splitlines():
            parsed_line = parse_line(line)
            if parsed_line:
                parsed_data.append(parsed_line)

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

    return {"records": output_records} and print(output_records)
