import peewee


class Secret(peewee.Model):
    hash = peewee.CharField()
    secret = peewee.TextField()
    created_at = peewee.DateTimeField()
    expires_at = peewee.DateTimeField()
    remaining_views = peewee.IntegerField()

    class Meta:
        database = peewee.SqliteDatabase("task_db.db")
