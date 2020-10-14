import peewee
from models import Secret
import datetime
import hashlib


class DbHandler:

    @staticmethod
    def initialize():
        db = peewee.SqliteDatabase("task_db.db")
        db.connect()
        db.create_tables([Secret])
        db.close()

    @staticmethod
    def get_secret(hash_to_find):
        db = peewee.SqliteDatabase("task_db.db")
        db.connect()
        query = Secret.select().where(Secret.hash == hash_to_find)
        current_time = datetime.datetime.today()
        if query.scalar() is None:
            db.close()
            return False
        else:
            query[0].remaining_views -= 1
            query[0].save()
            if query[0].remaining_views == -1:
                return False
            elif query[0].remaining_views == 0:
                new_query = Secret.delete().where(Secret.hash == hash_to_find)
                new_query.execute()
                db.close()
            # probably secrets after expiration date should by deleted automatically,
            # but I decided to delete them on request
            elif query[0].expires_at != query[0].created_at and (query[0].expires_at-current_time).total_seconds() < 0:
                new_query = Secret.delete().where(Secret.hash == hash_to_find)
                new_query.execute()
                db.close()
                return False
            return query[0]

    @staticmethod
    def post_secret(secret, availability, views):
        db = peewee.SqliteDatabase("task_db.db")
        db.connect()
        offset = datetime.timedelta(minutes=availability)
        current_time = datetime.datetime.today()
        available_till = current_time + offset
        # added timestamp to hash generation to prevent same hash for same secret text
        new_secret = bytes((secret+str(current_time)).encode())
        new_hash = hashlib.pbkdf2_hmac('sha256', new_secret, b'salt', 100000).hex()
        new_record = Secret.create(hash=new_hash,
                      secret=secret,
                      created_at=current_time,
                      expires_at=available_till,
                      remaining_views=views)
        db.close()
        return new_record
