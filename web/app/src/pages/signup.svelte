<script>
  import { Link } from "svelte-routing";
  import { Card, Button, Label, Input, Checkbox } from 'flowbite-svelte';
  import HoneyLemonLogo from '../assets/BankLogo.png';
</script>

<Card>
  <form class="flex flex-col space-y-6" action="/">
    <div class="flex items-center justify-between">
      <h3 class="text-4xl font-bold text-lime-900 dark:text-white">Signup</h3>
      <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
    <Label class="space-y-2">
      <span class="text-gray-400">Idcard</span>
      <Input type="text" name="idcard" placeholder="13-Digit" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Full Name</span>
      <Input type="text" name="fullname" placeholder="John Doe" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Birthdate</span>
      <Input type="date" name="birthdate" placeholder="yyyy-mm-dd" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Address</span>
      <Input type="text" name="address" placeholder="Address" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Email</span>
      <Input type="email" name="email" placeholder="name@company.com" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Password</span>
      <Input type="password" name="password" placeholder="••••••••" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Confirm password</span>
      <Input type="password" name="confirmpassword" placeholder="••••••••" required />
    </Label>
    <Link to="/">
      <Button type="submit" class="w-full bg-green-400 hover:bg-green-500">Signup</Button>
    </Link>
    <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
      Already registered? <a href="/login" class="text-green-400 hover:text-green-500 dark:text-primary-500">Login</a>
    </div>
  </form>
</Card>
