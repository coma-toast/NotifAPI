import * as fs from 'fs'
import * as yaml from 'js-yaml'

const configFile = fs.readFileSync("../config.yaml")
const config = yaml.load(configFile)

// @ts-ignore
const beamsClient = new PusherPushNotifications.Client({
    instanceId: config.instanceid
});

beamsClient
    .start()
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
