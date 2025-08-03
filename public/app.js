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
      `[data-id="generate-error-message"]`
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

    function validateUrl(url) {
      try {
        const urlObj = new URL(url);
        return urlObj.protocol === 'http:' || urlObj.protocol === 'https:';
      } catch {
        return false;
      }
    }

    function showError(message) {
      $errorMessage.innerText = message;
      $errorMessage.classList.remove("opaque");
      setTimeout(() => {
        $errorMessage.classList.add("opaque");
      }, 4000);
    }

    function resetButton() {
      $submitButton.innerText = "Generate";
      $submitButton.classList.remove("error", "done", "loading");
      $submitButton.disabled = false;
    }

    $submitButton.addEventListener("click", () => {
      const link = $linkInput.value.trim();

      if (!link) {
        showError("Please enter a URL to shorten");
        $linkInput.focus();
        return;
      }

      if (!validateUrl(link)) {
        showError("Please enter a valid URL (must start with http:// or https://)");
        $linkInput.focus();
        return;
      }

      $submitButton.innerText = "Generating...";
      $submitButton.disabled = true;
      $submitButton.classList.add("loading");

      fetch("/api/link", {
        method: "POST",
        body: JSON.stringify({ link }),
      })
        .then((res) => {
          if (res.ok) {
            return res.json();
          }

          res.json().then(({ message }) => {
            showError(message || "Failed to shorten URL. Please try again.");
            $submitButton.innerText = "Error";
            $submitButton.classList.remove("loading");
            $submitButton.classList.add("error");
            setTimeout(resetButton, 3000);
          }).catch(() => {
            showError("Failed to shorten URL. Please try again.");
            $submitButton.innerText = "Error";
            $submitButton.classList.remove("loading");
            $submitButton.classList.add("error");
            setTimeout(resetButton, 3000);
          });
        })
        .then((data) => {
          if (!data || !data.hash) {
            showError("Failed to generate short URL. Please try again.");
            $submitButton.innerText = "Error";
            $submitButton.classList.remove("loading");
            $submitButton.classList.add("error");
            setTimeout(resetButton, 3000);
            return;
          }

          const shortUrl = href + data.hash;
          
          $noLinkText.classList.add("hidden");
          $showLink.classList.remove("hidden");
          $showLink.innerText = shortUrl;

          $submitButton.innerText = "✓ Generated";
          $submitButton.classList.remove("loading");
          $submitButton.classList.add("done");
          
          // Auto-focus the result for better UX
          setTimeout(() => {
            $showLink.scrollIntoView({ behavior: 'smooth', block: 'center' });
          }, 100);
        }).catch(() => {
          showError("Failed to generate short URL. Please try again.");
          $submitButton.innerText = "Error";
          $submitButton.classList.remove("loading");
          $submitButton.classList.add("error");
          setTimeout(resetButton, 3000);
        });
    });

    function copyToClipboard() {
      const url = $showLink.innerText;
      
      if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(url).then(() => {
          $copiedToClipboard.classList.remove("opaque");
          setTimeout(() => {
            $copiedToClipboard.classList.add("opaque");
          }, 3000);
        }).catch(() => {
          fallbackCopyToClipboard(url);
        });
      } else {
        fallbackCopyToClipboard(url);
      }
    }

    $showLink.addEventListener("click", copyToClipboard);
    
    $showLink.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        copyToClipboard();
      }
    });
    
    function fallbackCopyToClipboard(text) {
      const textArea = document.createElement("textarea");
      textArea.value = text;
      textArea.style.position = "fixed";
      textArea.style.left = "-999999px";
      textArea.style.top = "-999999px";
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      
      try {
        document.execCommand('copy');
        $copiedToClipboard.classList.remove("opaque");
        setTimeout(() => {
          $copiedToClipboard.classList.add("opaque");
        }, 3000);
      } catch (err) {
        console.error('Fallback copy failed:', err);
      }
      
      document.body.removeChild(textArea);
    }
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
