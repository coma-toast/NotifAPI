  const beamsClient = new PusherPushNotifications.Client({
    instanceId: '6e482588-a9a1-45a9-b786-2d367fc69eef',
  });

  beamsClient.start()
    .then(() => beamsClient.addDeviceInterest('hello'))
    .then(() => console.log('Successfully registered and subscribed!'))
    .catch(console.error);
