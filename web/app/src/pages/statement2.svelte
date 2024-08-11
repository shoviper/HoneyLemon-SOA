<script>
    import { onMount } from 'svelte';
    import { Card } from 'flowbite-svelte';
    import { Link } from 'svelte-routing';
    import BackButton from '../assets/back.png';
    
    //test case
    let clientName = "Abdul";
    let clientID = "123456780";
    let accountID = "2806152770";
    let start = "2024-08-07T00:00:00Z";
    let end = "2024-08-10T23:59:59Z";
  
    let activity = [
      {
        "txID": "3",
        "from": "2806152770",
        "to": "5157693998",
        "amount": 20,
        "type": "payment",
        "timestamp": "2024-08-07T11:28:58.18205+07:00"
      },
      {
        "txID": "6",
        "from": "2806152770",
        "to": "5157693998",
        "amount": 20,
        "type": "transaction",
        "timestamp": "2024-08-09T11:28:52.035985+07:00"
      },
      {
        "txID": "1",
        "from": "2806152770",
        "to": "5157693998",
        "amount": 20,
        "type": "payment",
        "timestamp": "2024-08-09T23:59:44.223252+07:00"
      }
    ];
  
    let filteredActivity = [];
  
    onMount(() => {
      filteredActivity = activity.filter(
        (tx) => new Date(tx.timestamp) >= new Date(start) && new Date(tx.timestamp) <= new Date(end)
      );
    });
  </script>
  
  <div class="flex flex-col items-center">
    <Card class="bg-white mb-5" size="lg" padding="xl" style="width: 900px;">
      <div class="flex flex-col h-full">
        <div class="flex items-center justify-between mb-4">
          <Link to="/statement">
            <img src="{BackButton}" class="h-4 w-4" alt="Back" />
          </Link>
        </div>
        <div class="mb-4 text-gray-700 dark:text-gray-400">
            <h5 class="text-[#004D00] mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
            Statement for {clientName}
            </h5>
            <p><strong>Account ID:</strong> {accountID}</p>
            <p><strong>Client ID:</strong> {clientID}</p>
            <p><strong>Period:</strong> {new Date(start).toLocaleDateString()} to {new Date(end).toLocaleDateString()}</p>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full text-left text-sm">
            <thead>
              <tr class="border-b border-gray-200">
                <th class="py-2">Transaction ID</th>
                <th class="py-2">From</th>
                <th class="py-2">To</th>
                <th class="py-2">Amount</th>
                <th class="py-2">Type</th>
                <th class="py-2">Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {#each filteredActivity as tx}
                <tr class="border-b border-gray-200">
                  <td class="py-2">{tx.txID}</td>
                  <td class="py-2">{tx.from}</td>
                  <td class="py-2">{tx.to}</td>
                  <td class="py-2">${tx.amount}</td>
                  <td class="py-2 capitalize">{tx.type}</td>
                  <td class="py-2">{new Date(tx.timestamp).toLocaleString()}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </Card>
  </div>
  