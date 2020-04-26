# Mauritius Sea Cables (Go)
Go server for Mauritius Sea Cables

---

## Using bash to upload `./data/realtime.json` content to a github gist

1. Create a gist on github.
2. Create a github token with scope `gists - create gists`.
3. Edit `update_gist` bash script and change `GIST_ID` and `GITHUB_TOKEN` approriately.
4. Make the bash script executable.
    ```shell
    $ chmod +x update_gist
    ```

5. Execute the script and enjoy!
    ```shell
    $ ./update_gist
    ```
