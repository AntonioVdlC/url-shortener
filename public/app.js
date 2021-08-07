window.addEventListener("DOMContentLoaded", ignite);

function ignite() {
  const pathname = window.location.pathname;
  const href = window.location.href;
  const page = pathname === "/" ? "generate" : "redirect";

  document.querySelector(`[data-page=${page}]`).classList.remove("hidden");

  // Generate
  if (page === "generate") {
    const $linkInput = document.querySelector(`[data-id="input-link"]`);
    const $submitButton = document.querySelector(`[data-id="submit-button"]`);

    const $errorMessage = document.querySelector(
      `[data-id="generate-error-message]`
    );

    const $noLinkText = document.querySelector(`[data-id="no-link-text"]`);
    const $showLink = document.querySelector(`[data-id="show-link"]`);

    const $copiedToClipboard = document.querySelector(
      `[data-id="copied-to-clipboard"]`
    );

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
    const $fullLink = document.querySelector(`[data-id="show-full-link"]`);

    const $cancelButton = document.querySelector(`[data-id="cancel-button"]`);
    const $continueButton = document.querySelector(
      `[data-id="continue-button"]`
    );

    const $errorMessage = document.querySelector(
      `[data-id="redirect-error-message"]`
    );

    const hash = pathname.slice(1);
    if (!hash) {
      return;
    }


    $cancelButton.addEventListener("click", () => {
      // Note: we cannot close the current window, so the next best thing is to go back
      window.history.go(-1);
    });

    fetch(`/api/link?hash=${hash}`, { method: "GET" })
      .then((res) => {
        if (res.ok) {
          return res.json();
        }

        $errorMessage.innerText = "Link not found.";
        $errorMessage.classList.remove("opaque");

        $cancelButton.innerText = "Go back";
        $continueButton.classList.add("hidden");
      })
      .then((data) => {
        const link = data.link;

        if (!link) {
          $errorMessage.innerText = "Link not found.";
          $errorMessage.classList.remove("opaque");

          $cancelButton.innerText = "Go back";
          $continueButton.classList.add("hidden");
          return;
        }

        $fullLink.innerText = link;

        $continueButton.addEventListener("click", () => {
          window.location.replace(link);
        });
      });
  }
  // ---
}
