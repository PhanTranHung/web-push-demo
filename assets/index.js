const GET_VAPID_PUBLIC_KEY = "/vapid/public-key";
const SEND_NOTIFICATION = "/notification";
const SAVE_SUBCRIPTION = "/notification/subcription";

{
  const isSupportSw = () => "serviceWorker" in navigator;
  const isSupportNoti = () => "Notification" in window;

  const checkServiceWorker = (htmlTag) => {
    const swTag = document.getElementById("feature-sw");
    if (isSupportSw()) swTag.textContent = "Service worker: supported";
    else swTag.textContent = "Service worker: not supported";

    const notiTag = document.getElementById("feature-noti");
    if (isSupportNoti()) notiTag.textContent = "Notification: supported";
    else notiTag.textContent = "Notification: not supported";

    const notiSubTag = document.getElementById("noti-sub");
    if (isSupportSw() && isSupportNoti()) {
      navigator.serviceWorker.ready
        .then((registration) => registration.pushManager.getSubscription())
        .then(
          (subscription) =>
            (notiSubTag.textContent = JSON.stringify(subscription))
        );
    } else {
      notiSubTag.textContent =
        "Require bold Service Worker and Notification are supported";
    }
  };
  checkServiceWorker();

  function registerServiceWorker() {
    return navigator.serviceWorker.register("/service-worker.js", {
      type: "classic",
    });
  }
  if ("serviceWorker" in navigator) {
    registerServiceWorker().then(console.log);
  }
}

{
  function subscribe(vapidPublicKey) {
    return navigator.serviceWorker.ready.then(function (registration) {
      return registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
      });
    });
  }

  function urlBase64ToUint8Array(str) {
    const padding = "=".repeat((4 - (str.length % 4)) % 4);
    const base64 = (str + padding).replace(/\-/g, "+").replace(/_/g, "/");

    console.log(base64);
    const rawData = window.atob(base64);
    return Uint8Array.from([...rawData].map((char) => char.charCodeAt(0)));
  }

  function askNotificationPermission() {
    if (!("serviceWorker" in navigator) || !("Notification" in window)) {
      throw new Error("Browser is not support");
    }
    return new Promise((resolve, reject) => {
      const permissionResult = Notification.requestPermission((result) => {
        resolve(result);
      });

      if (permissionResult) {
        permissionResult.then(resolve, reject);
      }
    });
  }

  function subscribeNotification() {
    askNotificationPermission()
      .then((permissionResult) => {
        if (permissionResult !== "granted")
          throw new Error("Permission is not granted");
        return navigator.serviceWorker.ready;
      })
      .then((registration) => {
        return registration.pushManager.getSubscription();
      })
      .then((subscription) => {
        if (subscription) return subscription.unsubscribe();
        return true;
      })
      .then(() =>
        fetch(GET_VAPID_PUBLIC_KEY)
          .then((res) => res.text())
          .then((vapidPublicKey) => {
            console.log("vapidPublicKey", vapidPublicKey);
            return subscribe(vapidPublicKey);
          })
      )
      .then((subscription) => {
        const sub = JSON.stringify(subscription);
        document.getElementById("noti-sub").textContent = sub;
        console.log("subscription", sub);
        axios({
          method: "POST",
          url: SAVE_SUBCRIPTION,
          data: subscription,
        })
          .then((res) => console.log(res))
          .catch((err) => console.error(err));
      });
  }
}

{
  const getCustomerManifest = () => {
    return {
      name: "Web push",
      short_name: "Web push",
      description: "Demo of web push notification",
      theme_color: "#ffffff",
      start_url: window.location.origin,
      display: "standalone",
      background_color: "#ffffff",
      icons: [
        {
          src: window.location.origin + "/pos-192x192.png",
          sizes: "192x192",
          type: "image/png",
        },
        {
          src: window.location.origin + "/pos-512x512.png",
          sizes: "512x512",
          type: "image/png",
        },
        {
          src: window.location.origin + "/pos-512x512.png",
          sizes: "512x512",
          type: "image/png",
          purpose: "any maskable",
        },
      ],
    };
  };

  const attachManifestToDom = (manifest) => {
    const stringManifest = JSON.stringify(manifest);
    const blob = new Blob([stringManifest], { type: "application/json" });
    const manifestURL = URL.createObjectURL(blob);

    const manifestElement = document.createElement("link");
    manifestElement.setAttribute("rel", "manifest");
    manifestElement.setAttribute("href", manifestURL);
    document.head.appendChild(manifestElement);
  };

  const manifest = getCustomerManifest();
  console.log(manifest);
  attachManifestToDom(manifest);
}

{
  document.getElementById("noti-form").addEventListener("submit", (e) => {
    e.preventDefault();
    console.log(e);
    const f = new FormData(e.target);
    const payload = f.get("payload");
    axios({
      method: "POST",
      url: SEND_NOTIFICATION,
      data: JSON.parse(payload),
    })
      .then((res) => {
        document.getElementById("noti-res").textContent = res.data;
        console.log(res);
      })
      .catch((err) => console.error(err));
  });
}
