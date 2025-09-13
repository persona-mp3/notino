const version ='v2.3';

self.addEventListener('install', () => {
    console.log('Installing...')
});

self.addEventListener('activate',() => {
    console.log('SW Activated');

    caches.keys()
        .then((currentCaches) =>{
            
            const deletedPromise = currentCaches.map((existingCache) => {

                if (existingCache!== version){
                    console.log('Deleting Caches');
                    return caches.delete(existingCache);
                }
            });
            return Promise.all(deletedPromise)
        })
});

self.addEventListener('fetch', (e) => {
    e.respondWith(
        (async () => {
            // Handle non-GET requests
            if (e.request.method !== 'GET') {
                return fetch(e.request);
            }

            try {
                const res = await fetch(e.request);
                const clonedRes = res.clone();
                
                try {
                    // Open cache with current version
                    const newCache = await caches.open(version);
                    await newCache.put(e.request, clonedRes);
                } catch(err) {
                    console.error(`Error occurred while cloning: ${err}`);
                }
                
                return res;
            } catch(err) {
                console.error(`Error during request event`);
                // Fallback option for cached response
                const cache = await caches.open(version);
                const cachedResponse = await cache.match(e.request);
                if (cachedResponse) {
                    return cachedResponse;
                } else {
                    return caches.match('/index.html');
                }
            }
        })()
    );
});