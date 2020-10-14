from urllib.parse import parse_qs


class PostHandler:

    @staticmethod
    def check_validity(posted_data):
        parsed_to_dict = parse_qs(posted_data)
        try:
            if int(parsed_to_dict['expireAfterViews'][0]) <= 0 or len(parsed_to_dict["expireAfterViews"]) != 1:
                return False
            if int(parsed_to_dict['expireAfter'][0]) < 0 or len(parsed_to_dict["expireAfter"]) != 1:
                return False
            if "secret" not in parsed_to_dict.keys() or len(parsed_to_dict["secret"]) != 1:
                return False
            return True
        except KeyError:
            return False

    @staticmethod
    def unpack_values(posted_data):
        parsed_to_dict = parse_qs(posted_data)
        prepared_dict = {}
        for key in parsed_to_dict.keys():
            if parsed_to_dict[key][0].isnumeric():
                prepared_dict[key] = int(parsed_to_dict[key][0])
            else:
                prepared_dict[key] = parsed_to_dict[key][0]
        return prepared_dict
