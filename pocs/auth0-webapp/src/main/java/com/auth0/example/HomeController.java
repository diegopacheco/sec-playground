package com.auth0.example;

import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;

@Controller
public class HomeController {

	@GetMapping("/")
	String home(Model model, @AuthenticationPrincipal OidcUser user) {
		if (user != null) {
			return "redirect:/profile";
		}
		model.addAttribute("user", user);
		return "index";
	}

	@GetMapping("/profile")
	String profile(Model model, @AuthenticationPrincipal OidcUser user) {
		model.addAttribute("user", user);
		model.addAttribute("claims", user.getClaims());
		return "profile";
	}
}
