window.addEventListener("DOMContentLoaded", ignite);

function ignite() {
  const pathname = window.location.pathname;
  const href = window.location.href;
  const page = pathname === "/" ? "generate" : "redirect";

  document.querySelector(`[data-page=${page}]`).classList.remove("hidden");

  // Generate
  if (page === "generate") {
    const $linkInput = document.querySelector(`input[name="link"]`);
    const $submitButton = document.querySelector(".submit-button");

    const $errorMessage = document.querySelector(".error-message");

    const $noLinkText = document.querySelector(".no-link-text");
    const $showLink = document.querySelector(".show-link");

    const $copiedToClipboard = document.querySelector(".copied-to-clipboard");

    $linkInput.addEventListener("keypress", (e) => {
      if (e.key === "Enter") {
        $submitButton.click();
      }
    });

    $submitButton.addEventListener("click", () => {
      const link = $linkInput.value;

      if (!link) {
        // FIXME: error handling
        return;
      }

      $submitButton.innerText = "...";
      $submitButton.disabled = true;

      fetch("/api/link", {
        method: "POST",
        body: JSON.stringify({ link }),
      })
        .then((res) => {
          if (res.ok) {
            return res.json();
          }

          res.json().then(({ message }) => {
            $errorMessage.innerText = message || "Error";
            $errorMessage.classList.remove("opaque");

            $submitButton.innerText = "Error";
            $submitButton.classList.add("error");

            setTimeout(() => {
              $errorMessage.classList.add("opaque");

              $submitButton.innerText = "Generate";
              $submitButton.classList.remove("error");
              $submitButton.disabled = false;
            }, 2000);
          });
        })
        .then((data) => {
          const hash = data.hash;

          if (!hash) {
            $errorMessage.innerText = "Oops, an error has occured.";
            $errorMessage.classList.remove("opaque");

            $submitButton.innerText = "Error";
            $submitButton.classList.add("error");

            setTimeout(() => {
              $errorMessage.classList.add("opaque");

              $submitButton.innerText = "Generate";
              $submitButton.classList.remove("error");
              $submitButton.disabled = false;
            }, 2000);

            return;
          }

          $noLinkText.classList.add("hidden");
          $showLink.classList.remove("hidden");

          $showLink.innerText = href + hash;

          $submitButton.innerText = "Generated";
          $submitButton.classList.add("done");
        });
    });

    $showLink.addEventListener("click", () => {
      navigator.clipboard.writeText($showLink.innerText).then(() => {
        $copiedToClipboard.classList.remove("opaque");

        setTimeout(() => {
          $copiedToClipboard.classList.add("opaque");
        }, 2000);
      });
    });
  }
  // ---

  // Redirect
  if (page === "redirect") {
    const hash = pathname.slice(1);
    if (!hash) {
      // FIXME: error handling
      return;
    }

    // TODO: if request too slow, display waiting screen?

    fetch(`/api/link?hash=${hash}`, { method: "GET" })
      .then((res) => {
        // FIXME: error handling
        if (!res.ok) {
          throw new Error();
        }

        return res.json();
      })
      .then((data) => {
        const link = data.link;

        if (!link) {
          // FIXME: error handling
          return;
        }

        window.location.replace(link);
      });
  }
  // ---
}
