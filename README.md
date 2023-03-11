# Wg-weather

Redistributes requests from the https://weather-api.wundergraph.com/ GraphQL API as JSON.

The server instance exits after 10 seconds of inactivity.

# Challenges I had along the way


- While deploying to fly.io I tried to use Terraform, but I found that in order to deploy I would need to specify the container name, something that's not available through the fly_app Datasource.
  - Alternatively I could construct this out of `registry.fly.io/hyperized-weather:deployment-xxxx` but then I ran into the second bit:
  - I don't have a Terraform Cloud account, so I couldn't save the state to utilize Github Actions to do these deployments for me.
- When you request a flyctl allocate-v4 shared via the fly cli, it'll prompt you about incurring costs, even-though shared should be free.
- While deploying via the Fly.io Github actions I found myself unable to find a good way to ensure the app / machine state via the flyctl utility.
- Short on time, I opted for a quick naive implementation. Each push will simply remove and redeploy the application.