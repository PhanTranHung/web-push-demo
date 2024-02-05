/**
 *
 *
 *  ========================
 *
 *  DO NOT CHANGE THE NAME OF THIS FILE!
 *
 *  ========================
 *
 */

self.addEventListener("activate", () => self.clients.claim());
self.skipWaiting();

self.addEventListener("install", (evt) => {
  console.log("[Service Worker] Service worker installed");
});

self.addEventListener("fetch", function (event) {
  console.log("[Service Worker] Fetch event!\n", event);
});

self.addEventListener("push", (event) => {
  console.log("[Service Worker] Push Received.");
  console.log(`[Service Worker] Push had this data: "${event.data.text()}"`);

  const payload = JSON.parse(event.data.text());
  console.log(payload);

  event.waitUntil(
    self.registration.showNotification(payload.title, {
      body: payload.body,
      icon: payload.icon,
      vibrate: payload.vibrate,
    })
  );
});

self.addEventListener("notificationclick", (event) => {
  console.log(
    "[Service Worker] On notification click: ",
    event.notification.tag
  );
  event.notification.close();

  // This looks to see if the current is already open and
  // focuses if it is
  event.waitUntil(
    clients
      .matchAll({
        type: "window",
      })
      .then((clientList) => {
        for (const client of clientList) {
          if (client.url === "/" && "focus" in client) return client.focus();
        }
        if (clients.openWindow) return clients.openWindow("/");
      })
  );
});
