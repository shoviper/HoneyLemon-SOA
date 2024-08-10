<script>
  import { users, currentUser } from '../lib/userstore.js';
  import { navigate } from 'svelte-routing';
  import { Card, Button, Label, Input } from 'flowbite-svelte';
  import HoneyLemonLogo from "../assets/BankLogo.png";

  let user = null;
  let localUsers = [];

  users.subscribe(value => {
    localUsers = value;
  });

  currentUser.subscribe(value => {
    user = value;
  });

  function addNewAccount(event) {
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

    // Add new account
    user.accounts.push({ accountNumber: accnumber });

    // Update the user in the users array
    const updatedUsers = localUsers.map(u => u.idcard === user.idcard ? user : u);
    users.set(updatedUsers);

    // Update the currentUser store
    currentUser.set(user);

    alert("Account added successfully");

    // Redirect to the account page
    navigate('/mainaccount');
  }
</script>

<Card class="w-full max-w-md mx-auto">
  <form class="flex flex-col space-y-6" on:submit={addNewAccount}>
    <div class="flex items-center justify-between">
      <h3 class="text-4xl font-bold text-[#004D00] dark:text-white">Add Account</h3>
      <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
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
    <Button type="submit" class="w-full bg-[#28A745] hover:bg-[#03C04A]">Add account</Button>
  </form>
</Card>
