<script>
  import { Card, Button, Label, Input, Checkbox } from 'flowbite-svelte';
  import HoneyLemonLogo from '../assets/BankLogo.png';
</script>

<Card class="w-full max-w-lg mx-auto">
  <form class="flex flex-col space-y-6" action="/">
    <div class="flex items-center justify-between">
      <h3 class="text-4xl font-bold text-lime-900 dark:text-white">Register Account</h3>
      <img src="{HoneyLemonLogo}" class="h-18 w-28" alt="HoneyLemonLogo" />
    </div>
    <Label class="space-y-2">
      <span class="text-gray-400">Full Name</span>
      <Input type="text" name="fullname" placeholder="John Doe" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Birthdate</span>
      <Input type="date" name="birthdate" placeholder="yyyy-mm-dd" required />
    </Label>
    <div class="space-y-2">
      <span class="text-gray-400">Gender</span>
      <div class="flex items-center space-x-4">
        <Label class="inline-flex items-center">
          <input type="radio" name="gender" value="male" required class="appearance-none h-4 w-4 border border-gray-300 rounded-full checked:bg-green-600 checked:border-transparent focus:outline-none transition duration-200" />
          <span class="ml-2 text-gray-700 dark:text-gray-300">Male</span>
        </Label>
        <Label class="inline-flex items-center">
          <input type="radio" name="gender" value="female" required class="appearance-none h-4 w-4 border border-gray-300 rounded-full checked:bg-green-600 checked:border-transparent focus:outline-none transition duration-200" />
          <span class="ml-2 text-gray-700 dark:text-gray-300">Female</span>
        </Label>
        <Label class="inline-flex items-center">
          <input type="radio" name="gender" value="other" required class="appearance-none h-4 w-4 border border-gray-300 rounded-full checked:bg-green-600 checked:border-transparent focus:outline-none transition duration-200" />
          <span class="ml-2 text-gray-700 dark:text-gray-300">Others</span>
        </Label>
      </div>
    </div>
    <Label class="space-y-2">
      <span class="text-gray-400">Address</span>
      <Input type="text" name="address" placeholder="Address" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Email</span>
      <Input type="email" name="email" placeholder="name@company.com" required />
    </Label>
    <Label class="space-y-2">
      <span class="text-gray-400">Telephone</span>
      <Input type="text" name="tel" placeholder="xxx-xxx-xxxx" required />
    </Label>
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

