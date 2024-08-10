<script>
  import { users, currentUser } from '../lib/userstore.js';
  import { navigate } from 'svelte-routing';
  import { Card, Button, Label, Input } from 'flowbite-svelte';
  import HoneyLemonLogo from '../assets/BankLogo.png';

  let localUsers = [];

  users.subscribe(value => {
    localUsers = value;
  });

  function signup(event) {
    event.preventDefault();

    const idcard = event.target.idcard.value;
    const fullname = event.target.fullname.value;
    const birthdate = event.target.birthdate.value;
    const address = event.target.address.value;
    const email = event.target.email.value;
    const password = event.target.password.value;
    const confirmpassword = event.target.confirmpassword.value;

    // Validate ID card length
    if (idcard.length != 13) {
      alert("ID card is not correct");
      return;
    }

    // Validate password and confirm password match
    if (password !== confirmpassword) {
      alert("Passwords do not match");
      return;
    }

    // Check if ID card is already in use
    if (localUsers.some(user => user.idcard === idcard)) {
      alert("This ID card has already registered");
      return;
    }

    // Add new user
    const newUser = {
      idcard,
      fullname,
      birthdate,
      address,
      email,
      password,
      accounts: [] // Initialize with an empty accounts array
    };
    localUsers.push(newUser);

    // Update the users store
    users.set(localUsers);

    // Set the currentUser store
    currentUser.set(newUser);

    alert("Sign up successful");

    // Navigate to the appropriate page
    navigate(newUser.accounts.length > 0 ? '/' : '/accountregister');
  }
</script>

<Card>
  <form class="flex flex-col space-y-6" on:submit={signup}>
    <!-- Form fields -->
    <Label class="space-y-2">
      <span class="text-gray-400">Idcard</span>
      <Input type="text" name="idcard" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Full Name</span>
      <Input type="text" name="fullname" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Birthdate</span>
      <Input type="date" name="birthdate" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Address</span>
      <Input type="text" name="address" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Email</span>
      <Input type="email" name="email" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Password</span>
      <Input type="password" name="password" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Confirm password</span>
      <Input type="password" name="confirmpassword" required />
    </Label>
    <Button type="submit" class="w-full bg-[#28A745] hover:bg-[#28A745] hover:opacity-70">Signup</Button>
    <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
      Already registered? <a href="/login" class="text-green-400 hover:text-green-500 dark:text-primary-500">Login</a>
    </div>
  </form>
</Card>