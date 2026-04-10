# Plan

This is the plan for the microblog of snonux.foo

## Style

The style should be the same as in ./design.html

* also link to foo.zone, my normal blog, at the top of every page of the microblog.
* dont have an archive link like in the design.html example.
* the "transmit to nexus" should just link to mailto:paul@nospan.buetow.org
* every microblog entry should have date/timestamp linke in the design.html exampl
* dont have the likes (hearts), comments, and rebroadcast buttons per blog post.
* we want to be able to navigat between the blog posts with cursor up-down keys and jk (vi style up-down keys)
* we also want to navigate with the left-right cursor keys to previous, next page (if any) and also hl(vi style up-down keys)
* when navigating, we want to highlight the active entry (with a border or background), so it's clear which post is selected.and also play a sound.
* when hit enter on a blog post, open it in a larger view. when hit ESC there, return. 

## Input dir

I want to create a microblog. It should have an input directory, where i can put multiple source file formats:

* Plain .txt
* Markdown .md
* Images (.png/.jpg/.gif)
* Audio (.mp3 files)

And the blog should automatically generate a post out of it. Furthermore, it should also support a Markdown with a reference to an image file in the same directory.

Once we invoke the microblog generator, it should process all files in the input dir into blog assets.

### Plain text input

Just add this as a plain text as a blog post

### Markdown .md

Add this as a formatted HTML entry with styling etc. For this, we need to convert the Markdown to HTML.

### Images

Just have the image displayed as its own entry.

### Audio

Just have a playback button for the audio in the resulting entry

### After being processed

After an input was procssed, remove the files from the input dir.

## Output dir

The output dir should only contain static assets. That the directory to be published via a webbrowser.

We expect one directory per blog post. The microblog generator then combines all of them together into multiple pages.

So we need to keep the individual directories per blog post since the pages need to be re-generated according to the total blog post count and same for the atom.xml feed, so we need the directories as intermediate formats. We can link to images directly to them, though. So the output directory format will be like this:


./outdir/index.html # Generated main page
./outdir/pageN.html # Older pages
./ourdir/posts/YYYY-MM-DD-HHmmss/ # Microblog post asset(s)

Downscale any images to a maximum width of 1024px and compress to 80% JPEG quality

## Page size limit

Have max 42 entries on the single HTML page. Once more, allow paging (e.g. go to next 42 pages, etc).

* Next page (if any) is only on the bottom of the page.
* Previous page button (if any) is at the top of the page.

## atom.xml feed

We should have an atom.xml feed with the last 42 entries generated every run.

## Comprehensive testing

* All features should have integration tests.

## Technologies used

* Implemented in Google Go
* Follow everything from the auditing-code-quality skill
