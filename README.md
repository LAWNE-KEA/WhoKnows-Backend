# Flask Variations

# A note about merging

Note the no-merge-readme strategy in `.gitattributes`:

```plaintext
README.md merge=no-merge-readme
```

This is to prevent merge conflicts in the README.md file across branches. 

In the `.git/config` file, you can add the following configuration:

```plaintext
[merge "no-merge-readme"]
    name = "Do not merge README.md files"
    driver = true
```

You can now merge normally:

```bash
$ git merge <branch-name>
```

## How to get started

Each branch is a tutorial in a different topic based on the same Flask application as in the `main` branch. 

One way to follow along is by:

1. Forking the repository to your own account.

2. Cloning the repository to your local machine.

3. Checking out the branch you are interested in (e.g. `git checkout <branch-name>`).

4. Following the instructions in the README of the branch.

5. You can now push changes to your own repository. 

## Pull requests

If you have any suggestions or improvements to the tutorials, feel free to open a pull request.


