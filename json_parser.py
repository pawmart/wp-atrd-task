import json


class MyParser:

    @staticmethod
    def parse_to_json(query):
        dict_to_parse = {}
        dict_to_parse["hash"] = query.hash
        dict_to_parse["secretText"] = query.secret
        dict_to_parse["createdAt"] = str(query.created_at)
        if query.created_at == query.expires_at:
            dict_to_parse["expiresAt"] = "never"
        else:
            dict_to_parse["expiresAt"] = str(query.expires_at)
        dict_to_parse["remainingViews"] = query.remaining_views
        parsed_to_json = json.dumps(dict_to_parse)
        return parsed_to_json

    @staticmethod
    def parse_message_to_json(message):
        dict_to_parse = {}
        dict_to_parse['message'] = message
        parsed_to_json = json.dumps(dict_to_parse)
        return parsed_to_json
