import argparse
import json
import re

class Message(object):
    def __init__(self, name):
        self.params = {}
        self.name = name

    def type(self):
        return "".join(x[0].title() + x[1:]for x in self.name.split("_"))


class Scalar(object):
    def __init__(self, value):
        self.value = value

    def type(self):
        try:
            int(self.value)
            return 'int64'
        except ValueError:
            pass

        try:
            float(self.value)
            return 'float'
        except ValueError:
            pass

        if isinstance(self.value, str):
            return 'string'
        if isinstance(self.value, bool):
            return 'bool'
        if isinstance(self.value, int):
            return 'int64'
        if isinstance(self.value, float):
            return 'float'
        return 'Any'


class Repeated(object):
    def __init__(self, value):
        self.value = value

    def type(self):
        return f"repeated {self.value.type()}"


def value(ctx, record, name=None):
    if isinstance(record, list):
        msg = Message(name)
        for item in record:
            msg = value(ctx, item, name=name)

        if isinstance(msg, Message):
            ctx[name] = msg

        return Repeated(msg)

    if isinstance(record, dict):
        msg = Message(name)
        for key, v in record.items():
            field = value(ctx, v, name=key)
            msg.params[key] = field

            if isinstance(field, Message):
                ctx[key] = field

        return msg

    return Scalar(record)


def render(ctx):
    print('syntax = "proto3";')
    print('package grain.twitter;')
    print('option go_package = "twitterpb";')

    for name, message in ctx.items():
        print(f"message {message.type()} {{")
        for i, (name, field) in enumerate(message.params.items()):
            print(f"  {field.type()} {name} = {i+1};")
        print("}")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('input', type=argparse.FileType('r'))
    parser.add_argument('name', type=str)
    args = parser.parse_args()

    js = args.input.readlines()
    js[0] = "[{\n"
    record = json.loads("".join(js))

    ctx = {}
    value(ctx, record, name=f"{args.name}_entry")

    # l = Message(f"{args.name}_list")
    # l.params["items"] = Repeated(ctx[f"{args.name}_entry"])
    # ctx[f"{args.name}_list"] = l

    render(ctx)

if __name__ == "__main__":
    main()
