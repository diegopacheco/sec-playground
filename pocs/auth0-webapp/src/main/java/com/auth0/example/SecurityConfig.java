package com.auth0.example;

import java.io.IOException;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.logout.LogoutHandler;
import org.springframework.web.servlet.support.ServletUriComponentsBuilder;

@Configuration
public class SecurityConfig {

	@Value("${okta.oauth2.issuer}")
	private String issuer;

	@Value("${okta.oauth2.client-id}")
	private String clientId;

	@Bean
	SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
		http
			.authorizeHttpRequests(authorize -> authorize
				.requestMatchers("/", "/css/**").permitAll()
				.anyRequest().authenticated())
			.oauth2Login(oauth -> oauth.defaultSuccessUrl("/profile", true))
			.logout(logout -> logout.addLogoutHandler(auth0LogoutHandler()));
		return http.build();
	}

	private LogoutHandler auth0LogoutHandler() {
		return (request, response, authentication) -> {
			try {
				String returnTo = ServletUriComponentsBuilder.fromCurrentContextPath().path("/").build().toUriString();
				response.sendRedirect(issuer + "v2/logout?client_id=" + clientId + "&returnTo=" + returnTo);
			}
			catch (IOException exception) {
				throw new IllegalStateException(exception);
			}
		};
	}
}
