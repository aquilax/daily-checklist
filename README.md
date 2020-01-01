# daily-checklist

Simple tool for generating daily checklists in markdown

```
NAME:
   daily-checklist - Daily checklist processor

USAGE:
   daily-checklist template_file_name [date]

   Processes the template and outputs the result to stdout. If date is not provided, current date will be used.

NOTES:
   Template file can be any text file. All lines will be printed to stdout by default.
   To control if a line should be outputted add "<!-- @ [dom] [mon] [dow] -->" segment to the line which will output it only if the date matches the template. Date matching uses cron formatting.
```

## Example usage

For the following example:

```markdown
## Daily checklist

* [ ] Check weight <!-- @ * * 0 -->
* [ ] Exercise <!-- @ * * 2,4 -->
* [ ] Work <!-- @ * * 1,2,3,4,5 -->
* [ ] Prepare personal budget for next month <!-- @ 25 * * -->
* [ ] Bought christmas presents <!-- @ 20 12 * -->
```

```markdown
daily-checklist examples/checklist.md 2000-01-01
## Daily checklist

```

```markdown
daily-checklist examples/checklist.md 2000-01-04
## Daily checklist

* [ ] Exercise <!-- @ * * 2,4 -->
* [ ] Work <!-- @ * * 1,2,3,4,5 -->
```

Resulting into:

## Daily checklist

* [ ] Exercise <!-- @ * * 2,4 -->
* [ ] Work <!-- @ * * 1,2,3,4,5 -->
