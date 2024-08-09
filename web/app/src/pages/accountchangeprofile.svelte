<script>
  import { Card, Button, Label, Input } from 'flowbite-svelte';
  import HoneyLemonLogo from '../assets/BankLogo.png';
</script>

<Card class="w-full max-w-lg mx-auto mt-20">
  <form class="flex flex-col space-y-6" action="/">
    <div class="flex items-center justify-between">
      <h3 class="text-4xl font-bold text-lime-900 dark:text-white">Edit Profile</h3>
      <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
    <Label class="space-y-2">
      <span class="text-gray-400">ID Card</span>
      <Input type="text" name="idcard" placeholder="JD510" value="1234567890123" required />
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
    <Button type="submit" class="w-full bg-green-400 hover:bg-green-500">Save Changes</Button>
  </form>
</Card>
