# Mauritius Sea Cables (Go)
Go server for Mauritius Sea Cables

---

## Using bash to upload `./data/realtime.json` content to a github gist

1. Create a gist on github.
2. Create a github token with scope `gists - create gists`.
3. Make a copy of `.env.dist` and rename it to `.env`.
    ```shell
    $ cp .env.dist .env
    ```
4. Edit `GIST_ID` and `GITHUB_TOKEN` in `.env` approriately.
5. Make the bash script executable.
    ```shell
    $ chmod +x update_gist
    ```

6. Execute the script and enjoy!
    ```shell
    $ ./update_gist
    ```
