<script>
  import { Card, Button, Label, Input } from 'flowbite-svelte';
  import HoneyLemonLogo from '../assets/BankLogo.png';
  import { currentUser, users } from '../lib/userstore.js';
  import { navigate } from 'svelte-routing';
  
  let user = null;
  let localUsers = [];

  // Subscribe to the current user and users store
  currentUser.subscribe(value => {
    user = value;
  });

  users.subscribe(value => {
    localUsers = value;
  });

  function changePin(event) {
    event.preventDefault();

    const password = event.target.password.value;
    const newPin = event.target.newpin.value;
    const confirmNewPin = event.target.confirmnewpin.value;

    // Check if password matches
    if (password !== user.password) {
      alert("Incorrect password");
      return;
    }

    // Check if new PIN matches current PIN
    if (newPin === user.pin) {
      alert("New PIN cannot be the same as the current PIN");
      return;
    }

    // Check if new PIN and confirm new PIN match
    if (newPin !== confirmNewPin) {
      alert("New PINs do not match");
      return;
    }

    // Update the user's PIN
    user.pin = newPin;

    // Update the users store
    const updatedUsers = localUsers.map(u => u.idcard === user.idcard ? user : u);
    users.set(updatedUsers);

    // Update the currentUser store
    currentUser.set(user);

    alert("PIN changed successfully");

    // Navigate to the main account page or wherever appropriate
    navigate('/mainaccount');
  }
</script>

<Card class="w-full max-w-md mx-auto">
  <form class="flex flex-col space-y-6" on:submit={changePin}>
    <div class="flex items-center justify-between">
      <h3 class="text-4xl font-bold text-lime-900 dark:text-white">Change Pin</h3>
      <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
    <Label class="space-y-2">
      <span class="text-gray-400">Enter your current password</span>
      <Input type="password" name="password" placeholder="Current password" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Set a new 6 digit pin</span>
      <Input type="password" name="newpin" placeholder="New PIN" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Confirm new pin</span>
      <Input type="password" name="confirmnewpin" placeholder="Confirm new PIN" required />
    </Label>
    <Button type="submit" class="w-full bg-green-400 hover:bg-green-500">Save Changes</Button>
  </form>
</Card>
