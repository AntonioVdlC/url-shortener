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

    const $noLinkText = document.querySelector(".no-link-text");
    const $showLink = document.querySelector(".show-link");

    $submitButton.addEventListener("click", () => {
      const link = $linkInput.value;

      if (!link) {
        // FIXME: error handling
        return;
      }

      fetch("/api/link", {
        method: "POST",
        body: JSON.stringify({ link }),
      })
        .then((res) => {
          // FIXME: error handling
          if (!res.ok) {
            throw new Error();
          }

          return res.json();
        })
        .then((data) => {
          const hash = data.hash;

          if (!hash) {
            // FIXME: error handling
            return;
          }

          $noLinkText.classList.add("hidden");
          $showLink.classList.remove("hidden");

          $showLink.innerText = href + hash;
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
