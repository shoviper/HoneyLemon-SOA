<script>
  import { users, currentUser } from '../lib/userstore.js';
  import { navigate } from "svelte-routing";
  import { Card, Button, Label, Input } from 'flowbite-svelte';
  import HoneyLemonLogo from '../assets/BankLogo.png';

  let localUsers = [];

  users.subscribe(value => {
    localUsers = value;
  });

  function login(event) {
    event.preventDefault();

    const idcard = event.target.idcard.value;
    const password = event.target.password.value;

    const user = localUsers.find(user => user.idcard === idcard && user.password === password);

    if (user) {
      currentUser.set(user);
      alert("Login successful");
      navigate('/');
    } else {
      alert("Invalid ID card or password");
    }
  }
</script>

<Card>
  <form class="flex flex-col space-y-6" on:submit={login}>
    <div class="flex items-center justify-between">
      <h3 class="text-4xl font-bold text-lime-900 dark:text-white">Login</h3>
      <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
    <Label class="space-y-2">
      <span class="text-gray-400">Idcard</span>
      <Input type="text" name="idcard" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Password</span>
      <Input type="password" name="password" required />
    </Label>
    <div class="flex items-start">
      <a href="/resetpassword" class="ml-auto text-sm text-green-400 hover:text-green-500 dark:text-primary-500">Forgot password?</a>
    </div>
    <Button type="submit" class="w-full bg-[#28A745] hover:bg-[#28A745] hover:opacity-70">Login</Button>
    <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
      Not registered? <a href="/signup" class="text-green-400 hover:text-green-500 dark:text-primary-500">Create account</a>
    </div>
  </form>
</Card>