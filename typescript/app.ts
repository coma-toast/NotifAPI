// @ts-ignore
const beamsClient = new PusherPushNotifications.Client({
    instanceId: "6e482588-a9a1-45a9-b786-2d367fc69eef"
});

beamsClient
    .start()
    // .then(() => beamsClient.addDeviceInterest('hello'))
    .then(() => console.log("Successfully registered and subscribed!"))
    .catch(console.error);

let list = document.getElementById("interest-list");
beamsClient.getDeviceInterests().then((interests) => {
    console.log(interests);
    interests.forEach((interest) => {
        const li = document.createElement("li");
        li.classList.add("list-group-item");
        li.textContent = interest;
        list.appendChild(li);
        console.log(list);
    });
});
