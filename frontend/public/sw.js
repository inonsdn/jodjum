self.addEventListener("push", (event) => {
    let payload = {
        title: "All-Life",
        body: "You have a new notification"
    }

    if (event.data) {
        try {
            payload = event.data.json()
        } catch {
            payload.body = event.data.text()
        }
    }

    event.waitUntil(
        self.registration.showNotification(payload.title, {
            body: payload.body,
            icon: "/icons/icon-192.png",
            badge: "/icons/badge-72.png",
            data: {
                url: payload.url || "/",
            },
        }),
    );
})

self.addEventListener("notificationclick", (event) => {
    event.notification.close();
    const url = event.notification.data?.url || "/";

    event.waitUntil(
        (async () => {
            const clientsList = await clients.matchAll({
                type: "window",
                includeUncontrolled: true,
            });

            // Focus an already-open tab on the same origin instead of opening a
            // new one, navigating it to the notification's target url.
            for (const client of clientsList) {
                if ("focus" in client) {
                    await client.focus();
                    if ("navigate" in client) {
                        await client.navigate(url);
                    }
                    return;
                }
            }

            await clients.openWindow(url);
        })(),
    );
});