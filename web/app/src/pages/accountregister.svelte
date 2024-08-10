<script>
  import { users, currentUser } from '../lib/userstore.js';
  import { navigate } from 'svelte-routing';
  import { Card, Button, Label, Input } from 'flowbite-svelte';

  let user = null;
  let localUsers = [];

  users.subscribe(value => {
    localUsers = value;
  });

  currentUser.subscribe(value => {
    user = value;
  });

  function registerAccount(event) {
    event.preventDefault();

    const accnumber = event.target.accnumber.value;
    const fileaccbook = event.target.fileaccbook.files[0];

    if (!accnumber || !fileaccbook) {
      alert("Please fill in all required fields.");
      return;
    }

    // Check if the user already has 5 accounts
    if (user.accounts.length >= 5) {
      alert("You cannot have more than 5 accounts.");
      return;
    }

    // Add new account with initial balance
    user.accounts.push({ accountNumber: accnumber, balance: 10000.00 });

    // Update the currentUser and users store
    const updatedUsers = localUsers.map(u => u.idcard === user.idcard ? user : u);
    users.set(updatedUsers);
    currentUser.set(user);

    alert("Account registration successful");

    // Navigate to set PIN page
    navigate('/accountsetpin');
  }
</script>

<Card class="w-full max-w-lg mx-auto">
  <form class="flex flex-col space-y-6" on:submit={registerAccount}>
    <Label class="space-y-2">
      <span class="text-gray-400">Account number</span>
      <Input type="text" name="accnumber" placeholder="Account number" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Picture of account book</span>
      <div class="flex items-center">
        <Input type="file" name="fileaccbook" required class="hidden" id="fileaccbook" />
        <label for="fileaccbook" class="cursor-pointer flex items-center justify-center w-full py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
          Choose file
        </label>
      </div>
    </Label>
    <Button type="submit" class="w-full bg-green-400 hover:bg-green-500">Register</Button>
  </form>
</Card>