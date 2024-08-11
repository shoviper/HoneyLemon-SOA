<script>
  import { users, currentUser } from '../lib/userstore.js';
  import { navigate } from 'svelte-routing';
  import { Card, Button, Label, Input, Select } from 'flowbite-svelte';

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

    const acctype = event.target.acctype.value;
    const pin = event.target.pin.value;
    const confirmpin = event.target.confirmpin.value;

    if (!acctype || !pin || !confirmpin) {
      alert("Please fill in all required fields.");
      return;
    }

    if (pin !== confirmpin) {
      alert("Pins do not match.");
      return;
    }

    user.accounts.push({ accountType: acctype, balance: 10000.00, pin });

    const updatedUsers = localUsers.map(u => u.idcard === user.idcard ? user : u);
    users.set(updatedUsers);
    currentUser.set(user);

    alert("Account registration successful");

    navigate('/mainaccount');
  }
</script>

<Card class="w-full max-w-lg mx-auto">
  <form class="flex flex-col space-y-6" on:submit={registerAccount}>
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
