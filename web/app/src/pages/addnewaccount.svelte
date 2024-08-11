<script>
  import { users, currentUser } from '../lib/userstore.js';
  import { navigate } from 'svelte-routing';
  import { Card, Button, Label, Input, Select } from 'flowbite-svelte';
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

<Card class="w-full max-w-lg mx-auto">
  <form class="flex flex-col space-y-6" on:submit={addNewAccount}>
    <Label class="space-y-2">
      <div class="text-xs text-black mb-2 mr-4 w-full">
        Type
        <Select name="acctype" class="w-full mt-1 text-xs mt-2">
            <option value="Saving">Saving</option>
            <option value="Credit">Credit</option>
            <option value="Interest">Interest</option>
            <option value="Loan">Loan</option>
        </Select>
      </div>
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Set a 6 digit pin</span>
      <Input type="password" name="pin" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Confirm Pin</span>
      <Input type="password" name="confirmpin" required />
    </Label>
    <Button type="submit" class="w-full bg-green-400 hover:bg-green-500">Register</Button>
  </form>
</Card>