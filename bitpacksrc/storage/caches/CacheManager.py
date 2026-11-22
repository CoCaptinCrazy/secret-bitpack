from bitpack.storage.caches.Cache import Cache
from bitpack.storage.caches.CacheDriver import CacheDriver


class CacheManager:

    def __init__(self):
        self.driver = CacheDriver()

    def create(self, name):
        return Cache(name)

    def save(self, cache):
        self.driver.save(cache)

    def load(self, name):
        return self.driver.load(name)

    def clear(self, name):
        cache = Cache(name)
        self.driver.save(cache)