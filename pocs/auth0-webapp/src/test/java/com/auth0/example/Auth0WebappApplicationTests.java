package com.auth0.example;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.security.test.web.servlet.request.SecurityMockMvcRequestPostProcessors.oidcLogin;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.redirectedUrl;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.view;

@SpringBootTest
@AutoConfigureMockMvc
class Auth0WebappApplicationTests {

	@Autowired
	private MockMvc mockMvc;

	@Test
	void homeIsPublic() throws Exception {
		mockMvc.perform(get("/"))
			.andExpect(status().isOk())
			.andExpect(view().name("index"));
	}

	@Test
	void profileRequiresLogin() throws Exception {
		mockMvc.perform(get("/profile"))
			.andExpect(status().is3xxRedirection());
	}

	@Test
	void authenticatedHomeRedirectsToProfile() throws Exception {
		mockMvc.perform(get("/").with(oidcLogin()))
			.andExpect(status().is3xxRedirection())
			.andExpect(redirectedUrl("/profile"));
	}

	@Test
	void profileLoadsAfterLogin() throws Exception {
		mockMvc.perform(get("/profile").with(oidcLogin()))
			.andExpect(status().isOk())
			.andExpect(view().name("profile"));
	}

}
