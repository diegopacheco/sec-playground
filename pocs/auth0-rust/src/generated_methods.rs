use crate::error::Auth0Error;
use crate::management::{ManagementApi, ManagementRequest, RawManagementApi};

impl ManagementApi {
    pub fn actions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.list")
    }
    pub fn actions_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.create")
    }
    pub fn actions_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.get")
    }
    pub fn actions_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.delete")
    }
    pub fn actions_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.update")
    }
    pub fn actions_deploy(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.deploy")
    }
    pub fn actions_test(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.test")
    }
    pub fn branding_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.get")
    }
    pub fn branding_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.update")
    }
    pub fn client_grants_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("client_grants.list")
    }
    pub fn client_grants_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("client_grants.create")
    }
    pub fn client_grants_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("client_grants.get")
    }
    pub fn client_grants_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("client_grants.delete")
    }
    pub fn client_grants_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("client_grants.update")
    }
    pub fn clients_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.list")
    }
    pub fn clients_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.create")
    }
    pub fn clients_preview_cimd_metadata(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.preview_cimd_metadata")
    }
    pub fn clients_register_cimd_client(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.register_cimd_client")
    }
    pub fn clients_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.get")
    }
    pub fn clients_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.delete")
    }
    pub fn clients_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.update")
    }
    pub fn clients_rotate_secret(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.rotate_secret")
    }
    pub fn connection_profiles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.list")
    }
    pub fn connection_profiles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.create")
    }
    pub fn connection_profiles_list_templates(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.list_templates")
    }
    pub fn connection_profiles_get_template(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.get_template")
    }
    pub fn connection_profiles_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.get")
    }
    pub fn connection_profiles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.delete")
    }
    pub fn connection_profiles_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connection_profiles.update")
    }
    pub fn connections_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.list")
    }
    pub fn connections_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.create")
    }
    pub fn connections_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.get")
    }
    pub fn connections_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.delete")
    }
    pub fn connections_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.update")
    }
    pub fn connections_check_status(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.check_status")
    }
    pub fn custom_domains_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.list")
    }
    pub fn custom_domains_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.create")
    }
    pub fn custom_domains_get_default(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.get_default")
    }
    pub fn custom_domains_set_default(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.set_default")
    }
    pub fn custom_domains_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.get")
    }
    pub fn custom_domains_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.delete")
    }
    pub fn custom_domains_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.update")
    }
    pub fn custom_domains_test(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.test")
    }
    pub fn custom_domains_verify(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("custom_domains.verify")
    }
    pub fn device_credentials_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("device_credentials.list")
    }
    pub fn device_credentials_create_public_key(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("device_credentials.create_public_key")
    }
    pub fn device_credentials_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("device_credentials.delete")
    }
    pub fn email_templates_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("email_templates.create")
    }
    pub fn email_templates_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("email_templates.get")
    }
    pub fn email_templates_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("email_templates.set")
    }
    pub fn email_templates_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("email_templates.update")
    }
    pub fn event_streams_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("event_streams.list")
    }
    pub fn event_streams_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("event_streams.create")
    }
    pub fn event_streams_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("event_streams.get")
    }
    pub fn event_streams_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("event_streams.delete")
    }
    pub fn event_streams_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("event_streams.update")
    }
    pub fn event_streams_test(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("event_streams.test")
    }
    pub fn events_subscribe(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("events.subscribe")
    }
    pub fn flows_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.list")
    }
    pub fn flows_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.create")
    }
    pub fn flows_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.get")
    }
    pub fn flows_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.delete")
    }
    pub fn flows_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.update")
    }
    pub fn forms_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("forms.list")
    }
    pub fn forms_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("forms.create")
    }
    pub fn forms_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("forms.get")
    }
    pub fn forms_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("forms.delete")
    }
    pub fn forms_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("forms.update")
    }
    pub fn groups_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.list")
    }
    pub fn groups_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.get")
    }
    pub fn groups_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.delete")
    }
    pub fn hooks_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.list")
    }
    pub fn hooks_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.create")
    }
    pub fn hooks_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.get")
    }
    pub fn hooks_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.delete")
    }
    pub fn hooks_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.update")
    }
    pub fn jobs_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.get")
    }
    pub fn log_streams_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("log_streams.list")
    }
    pub fn log_streams_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("log_streams.create")
    }
    pub fn log_streams_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("log_streams.get")
    }
    pub fn log_streams_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("log_streams.delete")
    }
    pub fn log_streams_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("log_streams.update")
    }
    pub fn logs_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("logs.list")
    }
    pub fn logs_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("logs.get")
    }
    pub fn network_acls_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("network_acls.list")
    }
    pub fn network_acls_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("network_acls.create")
    }
    pub fn network_acls_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("network_acls.get")
    }
    pub fn network_acls_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("network_acls.set")
    }
    pub fn network_acls_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("network_acls.delete")
    }
    pub fn network_acls_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("network_acls.update")
    }
    pub fn organizations_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.list")
    }
    pub fn organizations_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.create")
    }
    pub fn organizations_get_by_name(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.get_by_name")
    }
    pub fn organizations_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.get")
    }
    pub fn organizations_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.delete")
    }
    pub fn organizations_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.update")
    }
    pub fn prompts_get_settings(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.get_settings")
    }
    pub fn prompts_update_settings(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.update_settings")
    }
    pub fn rate_limit_policies_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rate_limit_policies.list")
    }
    pub fn rate_limit_policies_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rate_limit_policies.create")
    }
    pub fn rate_limit_policies_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rate_limit_policies.get")
    }
    pub fn rate_limit_policies_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rate_limit_policies.delete")
    }
    pub fn rate_limit_policies_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rate_limit_policies.update")
    }
    pub fn refresh_tokens_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("refresh_tokens.list")
    }
    pub fn refresh_tokens_revoke(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("refresh_tokens.revoke")
    }
    pub fn refresh_tokens_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("refresh_tokens.get")
    }
    pub fn refresh_tokens_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("refresh_tokens.delete")
    }
    pub fn refresh_tokens_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("refresh_tokens.update")
    }
    pub fn resource_servers_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("resource_servers.list")
    }
    pub fn resource_servers_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("resource_servers.create")
    }
    pub fn resource_servers_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("resource_servers.get")
    }
    pub fn resource_servers_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("resource_servers.delete")
    }
    pub fn resource_servers_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("resource_servers.update")
    }
    pub fn roles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.list")
    }
    pub fn roles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.create")
    }
    pub fn roles_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.get")
    }
    pub fn roles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.delete")
    }
    pub fn roles_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.update")
    }
    pub fn rules_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules.list")
    }
    pub fn rules_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules.create")
    }
    pub fn rules_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules.get")
    }
    pub fn rules_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules.delete")
    }
    pub fn rules_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules.update")
    }
    pub fn rules_configs_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules_configs.list")
    }
    pub fn rules_configs_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules_configs.set")
    }
    pub fn rules_configs_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("rules_configs.delete")
    }
    pub fn self_service_profiles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("self_service_profiles.list")
    }
    pub fn self_service_profiles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("self_service_profiles.create")
    }
    pub fn self_service_profiles_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("self_service_profiles.get")
    }
    pub fn self_service_profiles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("self_service_profiles.delete")
    }
    pub fn self_service_profiles_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("self_service_profiles.update")
    }
    pub fn sessions_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("sessions.get")
    }
    pub fn sessions_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("sessions.delete")
    }
    pub fn sessions_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("sessions.update")
    }
    pub fn sessions_revoke(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("sessions.revoke")
    }
    pub fn stats_get_active_users_count(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("stats.get_active_users_count")
    }
    pub fn stats_get_daily(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("stats.get_daily")
    }
    pub fn supplemental_signals_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("supplemental_signals.get")
    }
    pub fn supplemental_signals_patch(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("supplemental_signals.patch")
    }
    pub fn tickets_verify_email(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("tickets.verify_email")
    }
    pub fn tickets_change_password(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("tickets.change_password")
    }
    pub fn token_exchange_profiles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("token_exchange_profiles.list")
    }
    pub fn token_exchange_profiles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("token_exchange_profiles.create")
    }
    pub fn token_exchange_profiles_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("token_exchange_profiles.get")
    }
    pub fn token_exchange_profiles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("token_exchange_profiles.delete")
    }
    pub fn token_exchange_profiles_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("token_exchange_profiles.update")
    }
    pub fn user_attribute_profiles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.list")
    }
    pub fn user_attribute_profiles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.create")
    }
    pub fn user_attribute_profiles_list_templates(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.list_templates")
    }
    pub fn user_attribute_profiles_get_template(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.get_template")
    }
    pub fn user_attribute_profiles_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.get")
    }
    pub fn user_attribute_profiles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.delete")
    }
    pub fn user_attribute_profiles_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_attribute_profiles.update")
    }
    pub fn user_blocks_list_by_identifier(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_blocks.list_by_identifier")
    }
    pub fn user_blocks_delete_by_identifier(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_blocks.delete_by_identifier")
    }
    pub fn user_blocks_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_blocks.list")
    }
    pub fn user_blocks_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_blocks.delete")
    }
    pub fn user_grants_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_grants.list")
    }
    pub fn user_grants_delete_by_user_id(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_grants.delete_by_user_id")
    }
    pub fn user_grants_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("user_grants.delete")
    }
    pub fn users_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.list")
    }
    pub fn users_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.create")
    }
    pub fn users_list_users_by_email(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.list_users_by_email")
    }
    pub fn users_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.get")
    }
    pub fn users_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.delete")
    }
    pub fn users_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.update")
    }
    pub fn users_regenerate_recovery_code(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.regenerate_recovery_code")
    }
    pub fn users_revoke_access(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.revoke_access")
    }
    pub fn actions_executions_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.executions.get")
    }
    pub fn actions_modules_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.list")
    }
    pub fn actions_modules_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.create")
    }
    pub fn actions_modules_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.get")
    }
    pub fn actions_modules_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.delete")
    }
    pub fn actions_modules_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.update")
    }
    pub fn actions_modules_list_actions(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.list_actions")
    }
    pub fn actions_modules_rollback(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.rollback")
    }
    pub fn actions_triggers_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.triggers.list")
    }
    pub fn actions_versions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.versions.list")
    }
    pub fn actions_versions_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.versions.get")
    }
    pub fn actions_versions_deploy(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.versions.deploy")
    }
    pub fn actions_modules_versions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.versions.list")
    }
    pub fn actions_modules_versions_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.versions.create")
    }
    pub fn actions_modules_versions_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.modules.versions.get")
    }
    pub fn actions_triggers_bindings_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.triggers.bindings.list")
    }
    pub fn actions_triggers_bindings_update_many(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("actions.triggers.bindings.update_many")
    }
    pub fn anomaly_blocks_check_ip(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("anomaly.blocks.check_ip")
    }
    pub fn anomaly_blocks_unblock_ip(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("anomaly.blocks.unblock_ip")
    }
    pub fn attackprotection_bot_detection_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.bot_detection.get")
    }
    pub fn attackprotection_bot_detection_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.bot_detection.update")
    }
    pub fn attackprotection_breached_password_detection_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.breached_password_detection.get")
    }
    pub fn attackprotection_breached_password_detection_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.breached_password_detection.update")
    }
    pub fn attackprotection_brute_force_protection_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.brute_force_protection.get")
    }
    pub fn attackprotection_brute_force_protection_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.brute_force_protection.update")
    }
    pub fn attackprotection_captcha_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.captcha.get")
    }
    pub fn attackprotection_captcha_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.captcha.update")
    }
    pub fn attackprotection_phone_provider_protection_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.phone_provider_protection.get")
    }
    pub fn attackprotection_phone_provider_protection_patch(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.phone_provider_protection.patch")
    }
    pub fn attackprotection_suspicious_ip_throttling_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.suspicious_ip_throttling.get")
    }
    pub fn attackprotection_suspicious_ip_throttling_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("attackprotection.suspicious_ip_throttling.update")
    }
    pub fn branding_templates_get_universal_login(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.templates.get_universal_login")
    }
    pub fn branding_templates_update_universal_login(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.templates.update_universal_login")
    }
    pub fn branding_templates_delete_universal_login(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.templates.delete_universal_login")
    }
    pub fn branding_themes_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.themes.create")
    }
    pub fn branding_themes_get_default(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.themes.get_default")
    }
    pub fn branding_themes_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.themes.get")
    }
    pub fn branding_themes_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.themes.delete")
    }
    pub fn branding_themes_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.themes.update")
    }
    pub fn branding_phone_providers_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.providers.list")
    }
    pub fn branding_phone_providers_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.providers.create")
    }
    pub fn branding_phone_providers_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.providers.get")
    }
    pub fn branding_phone_providers_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.providers.delete")
    }
    pub fn branding_phone_providers_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.providers.update")
    }
    pub fn branding_phone_providers_test(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.providers.test")
    }
    pub fn branding_phone_templates_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.list")
    }
    pub fn branding_phone_templates_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.create")
    }
    pub fn branding_phone_templates_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.get")
    }
    pub fn branding_phone_templates_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.delete")
    }
    pub fn branding_phone_templates_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.update")
    }
    pub fn branding_phone_templates_reset(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.reset")
    }
    pub fn branding_phone_templates_test(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("branding.phone.templates.test")
    }
    pub fn clientgrants_organizations_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clientgrants.organizations.list")
    }
    pub fn clients_connections_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.connections.get")
    }
    pub fn clients_credentials_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.credentials.list")
    }
    pub fn clients_credentials_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.credentials.create")
    }
    pub fn clients_credentials_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.credentials.get")
    }
    pub fn clients_credentials_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.credentials.delete")
    }
    pub fn clients_credentials_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("clients.credentials.update")
    }
    pub fn connections_clients_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.clients.get")
    }
    pub fn connections_clients_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.clients.update")
    }
    pub fn connections_directory_provisioning_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.list")
    }
    pub fn connections_directory_provisioning_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.get")
    }
    pub fn connections_directory_provisioning_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.create")
    }
    pub fn connections_directory_provisioning_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.delete")
    }
    pub fn connections_directory_provisioning_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.update")
    }
    pub fn connections_directory_provisioning_get_default_mapping(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.get_default_mapping")
    }
    pub fn connections_directory_provisioning_list_synchronized_groups(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.list_synchronized_groups")
    }
    pub fn connections_directory_provisioning_set(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directory_provisioning.set")
    }
    pub fn connections_keys_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.keys.get")
    }
    pub fn connections_keys_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.keys.create")
    }
    pub fn connections_keys_rotate(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.keys.rotate")
    }
    pub fn connections_scim_configuration_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scim_configuration.list")
    }
    pub fn connections_scim_configuration_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scim_configuration.get")
    }
    pub fn connections_scim_configuration_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scim_configuration.create")
    }
    pub fn connections_scim_configuration_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scim_configuration.delete")
    }
    pub fn connections_scim_configuration_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scim_configuration.update")
    }
    pub fn connections_scim_configuration_get_default_mapping(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scim_configuration.get_default_mapping")
    }
    pub fn connections_users_delete_by_email(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.users.delete_by_email")
    }
    pub fn connections_directoryprovisioning_synchronizations_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.directoryprovisioning.synchronizations.create")
    }
    pub fn connections_scimconfiguration_tokens_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scimconfiguration.tokens.get")
    }
    pub fn connections_scimconfiguration_tokens_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scimconfiguration.tokens.create")
    }
    pub fn connections_scimconfiguration_tokens_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("connections.scimconfiguration.tokens.delete")
    }
    pub fn emails_provider_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("emails.provider.get")
    }
    pub fn emails_provider_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("emails.provider.create")
    }
    pub fn emails_provider_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("emails.provider.delete")
    }
    pub fn emails_provider_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("emails.provider.update")
    }
    pub fn eventstreams_deliveries_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("eventstreams.deliveries.list")
    }
    pub fn eventstreams_deliveries_get_history(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("eventstreams.deliveries.get_history")
    }
    pub fn eventstreams_redeliveries_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("eventstreams.redeliveries.create")
    }
    pub fn eventstreams_redeliveries_create_by_id(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("eventstreams.redeliveries.create_by_id")
    }
    pub fn flows_executions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.executions.list")
    }
    pub fn flows_executions_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.executions.get")
    }
    pub fn flows_executions_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.executions.delete")
    }
    pub fn flows_vault_connections_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.vault.connections.list")
    }
    pub fn flows_vault_connections_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.vault.connections.create")
    }
    pub fn flows_vault_connections_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.vault.connections.get")
    }
    pub fn flows_vault_connections_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.vault.connections.delete")
    }
    pub fn flows_vault_connections_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("flows.vault.connections.update")
    }
    pub fn groups_members_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.members.get")
    }
    pub fn groups_roles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.roles.list")
    }
    pub fn groups_roles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.roles.create")
    }
    pub fn groups_roles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("groups.roles.delete")
    }
    pub fn guardian_enrollments_create_ticket(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.enrollments.create_ticket")
    }
    pub fn guardian_enrollments_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.enrollments.get")
    }
    pub fn guardian_enrollments_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.enrollments.delete")
    }
    pub fn guardian_factors_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.list")
    }
    pub fn guardian_factors_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.set")
    }
    pub fn guardian_policies_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.policies.list")
    }
    pub fn guardian_policies_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.policies.set")
    }
    pub fn guardian_factors_phone_get_message_types(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.get_message_types")
    }
    pub fn guardian_factors_phone_set_message_types(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.set_message_types")
    }
    pub fn guardian_factors_phone_get_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.get_twilio_provider")
    }
    pub fn guardian_factors_phone_set_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.set_twilio_provider")
    }
    pub fn guardian_factors_phone_get_selected_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.get_selected_provider")
    }
    pub fn guardian_factors_phone_set_provider(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.set_provider")
    }
    pub fn guardian_factors_phone_get_templates(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.get_templates")
    }
    pub fn guardian_factors_phone_set_templates(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.phone.set_templates")
    }
    pub fn guardian_factors_push_notification_get_apns_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.get_apns_provider")
    }
    pub fn guardian_factors_push_notification_set_apns_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_apns_provider")
    }
    pub fn guardian_factors_push_notification_update_apns_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_apns_provider")
    }
    pub fn guardian_factors_push_notification_set_fcm_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_fcm_provider")
    }
    pub fn guardian_factors_push_notification_update_fcm_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_fcm_provider")
    }
    pub fn guardian_factors_push_notification_set_fcmv1_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_fcmv1_provider")
    }
    pub fn guardian_factors_push_notification_update_fcmv1_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_fcmv1_provider")
    }
    pub fn guardian_factors_push_notification_get_sns_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.get_sns_provider")
    }
    pub fn guardian_factors_push_notification_set_sns_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_sns_provider")
    }
    pub fn guardian_factors_push_notification_update_sns_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_sns_provider")
    }
    pub fn guardian_factors_push_notification_get_selected_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.get_selected_provider")
    }
    pub fn guardian_factors_push_notification_set_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_provider")
    }
    pub fn guardian_factors_sms_get_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.sms.get_twilio_provider")
    }
    pub fn guardian_factors_sms_set_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.sms.set_twilio_provider")
    }
    pub fn guardian_factors_sms_get_selected_provider(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.sms.get_selected_provider")
    }
    pub fn guardian_factors_sms_set_provider(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.sms.set_provider")
    }
    pub fn guardian_factors_sms_get_templates(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.sms.get_templates")
    }
    pub fn guardian_factors_sms_set_templates(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.sms.set_templates")
    }
    pub fn guardian_factors_duo_settings_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.duo.settings.get")
    }
    pub fn guardian_factors_duo_settings_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.duo.settings.set")
    }
    pub fn guardian_factors_duo_settings_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("guardian.factors.duo.settings.update")
    }
    pub fn hooks_secrets_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.secrets.get")
    }
    pub fn hooks_secrets_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.secrets.create")
    }
    pub fn hooks_secrets_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.secrets.delete")
    }
    pub fn hooks_secrets_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("hooks.secrets.update")
    }
    pub fn jobs_errors_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.errors.get")
    }
    pub fn jobs_users_exports_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.users_exports.create")
    }
    pub fn jobs_users_imports_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.users_imports.create")
    }
    pub fn jobs_users_imports_create_post(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.users_imports.create_post")
    }
    pub fn jobs_users_imports_create_post_3(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.users_imports.create_post_3")
    }
    pub fn jobs_users_imports_create_post_4(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.users_imports.create_post_4")
    }
    pub fn jobs_users_imports_create_post_5(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.users_imports.create_post_5")
    }
    pub fn jobs_verification_email_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("jobs.verification_email.create")
    }
    pub fn keys_custom_signing_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.custom_signing.get")
    }
    pub fn keys_custom_signing_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.custom_signing.set")
    }
    pub fn keys_custom_signing_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.custom_signing.delete")
    }
    pub fn keys_encryption_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.list")
    }
    pub fn keys_encryption_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.create")
    }
    pub fn keys_encryption_rekey(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.rekey")
    }
    pub fn keys_encryption_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.get")
    }
    pub fn keys_encryption_import_(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.import_")
    }
    pub fn keys_encryption_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.delete")
    }
    pub fn keys_encryption_create_public_wrapping_key(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.encryption.create_public_wrapping_key")
    }
    pub fn keys_signing_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.signing.list")
    }
    pub fn keys_signing_rotate(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.signing.rotate")
    }
    pub fn keys_signing_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.signing.get")
    }
    pub fn keys_signing_revoke(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("keys.signing.revoke")
    }
    pub fn organizations_client_grants_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.client_grants.list")
    }
    pub fn organizations_client_grants_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.client_grants.create")
    }
    pub fn organizations_client_grants_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.client_grants.delete")
    }
    pub fn organizations_connections_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.connections.list")
    }
    pub fn organizations_connections_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.connections.create")
    }
    pub fn organizations_connections_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.connections.get")
    }
    pub fn organizations_connections_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.connections.delete")
    }
    pub fn organizations_connections_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.connections.update")
    }
    pub fn organizations_discovery_domains_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.discovery_domains.list")
    }
    pub fn organizations_discovery_domains_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.discovery_domains.create")
    }
    pub fn organizations_discovery_domains_get_by_name(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.discovery_domains.get_by_name")
    }
    pub fn organizations_discovery_domains_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.discovery_domains.get")
    }
    pub fn organizations_discovery_domains_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.discovery_domains.delete")
    }
    pub fn organizations_discovery_domains_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.discovery_domains.update")
    }
    pub fn organizations_enabled_connections_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.enabled_connections.list")
    }
    pub fn organizations_enabled_connections_add(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.enabled_connections.add")
    }
    pub fn organizations_enabled_connections_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.enabled_connections.get")
    }
    pub fn organizations_enabled_connections_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.enabled_connections.delete")
    }
    pub fn organizations_enabled_connections_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.enabled_connections.update")
    }
    pub fn organizations_groups_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.groups.list")
    }
    pub fn organizations_invitations_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.invitations.list")
    }
    pub fn organizations_invitations_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.invitations.create")
    }
    pub fn organizations_invitations_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.invitations.get")
    }
    pub fn organizations_invitations_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.invitations.delete")
    }
    pub fn organizations_members_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.list")
    }
    pub fn organizations_members_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.create")
    }
    pub fn organizations_members_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.delete")
    }
    pub fn organizations_groups_roles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.groups.roles.list")
    }
    pub fn organizations_groups_roles_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.groups.roles.create")
    }
    pub fn organizations_groups_roles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.groups.roles.delete")
    }
    pub fn organizations_members_effective_roles_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.effective_roles.list")
    }
    pub fn organizations_members_roles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.roles.list")
    }
    pub fn organizations_members_roles_assign(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.roles.assign")
    }
    pub fn organizations_members_roles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.roles.delete")
    }
    pub fn organizations_members_effectiveroles_sources_groups_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("organizations.members.effectiveroles.sources.groups.list")
    }
    pub fn prompts_custom_text_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.custom_text.get")
    }
    pub fn prompts_custom_text_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.custom_text.set")
    }
    pub fn prompts_partials_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.partials.get")
    }
    pub fn prompts_partials_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.partials.set")
    }
    pub fn prompts_rendering_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.rendering.list")
    }
    pub fn prompts_rendering_bulk_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.rendering.bulk_update")
    }
    pub fn prompts_rendering_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.rendering.get")
    }
    pub fn prompts_rendering_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("prompts.rendering.update")
    }
    pub fn riskassessments_settings_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("riskassessments.settings.get")
    }
    pub fn riskassessments_settings_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("riskassessments.settings.update")
    }
    pub fn riskassessments_settings_new_device_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("riskassessments.settings.new_device.get")
    }
    pub fn riskassessments_settings_new_device_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("riskassessments.settings.new_device.update")
    }
    pub fn roles_groups_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.groups.get")
    }
    pub fn roles_groups_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.groups.create")
    }
    pub fn roles_groups_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.groups.delete")
    }
    pub fn roles_permissions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.permissions.list")
    }
    pub fn roles_permissions_add(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.permissions.add")
    }
    pub fn roles_permissions_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.permissions.delete")
    }
    pub fn roles_users_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.users.list")
    }
    pub fn roles_users_assign(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("roles.users.assign")
    }
    pub fn selfserviceprofiles_custom_text_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("selfserviceprofiles.custom_text.list")
    }
    pub fn selfserviceprofiles_custom_text_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("selfserviceprofiles.custom_text.set")
    }
    pub fn selfserviceprofiles_sso_ticket_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("selfserviceprofiles.sso_ticket.create")
    }
    pub fn selfserviceprofiles_sso_ticket_revoke(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("selfserviceprofiles.sso_ticket.revoke")
    }
    pub fn tenants_settings_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("tenants.settings.get")
    }
    pub fn tenants_settings_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("tenants.settings.update")
    }
    pub fn users_authentication_methods_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.list")
    }
    pub fn users_authentication_methods_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.create")
    }
    pub fn users_authentication_methods_set(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.set")
    }
    pub fn users_authentication_methods_delete_all(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.delete_all")
    }
    pub fn users_authentication_methods_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.get")
    }
    pub fn users_authentication_methods_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.delete")
    }
    pub fn users_authentication_methods_update(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authentication_methods.update")
    }
    pub fn users_authenticators_delete_all(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.authenticators.delete_all")
    }
    pub fn users_connected_accounts_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.connected_accounts.list")
    }
    pub fn users_effective_permissions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.effective_permissions.list")
    }
    pub fn users_effective_roles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.effective_roles.list")
    }
    pub fn users_enrollments_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.enrollments.get")
    }
    pub fn users_federated_connections_tokensets_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.federated_connections_tokensets.list")
    }
    pub fn users_federated_connections_tokensets_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.federated_connections_tokensets.delete")
    }
    pub fn users_groups_get(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.groups.get")
    }
    pub fn users_identities_link(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.identities.link")
    }
    pub fn users_identities_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.identities.delete")
    }
    pub fn users_logs_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.logs.list")
    }
    pub fn users_multifactor_invalidate_remember_browser(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.multifactor.invalidate_remember_browser")
    }
    pub fn users_multifactor_delete_provider(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.multifactor.delete_provider")
    }
    pub fn users_organizations_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.organizations.list")
    }
    pub fn users_permissions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.permissions.list")
    }
    pub fn users_permissions_create(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.permissions.create")
    }
    pub fn users_permissions_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.permissions.delete")
    }
    pub fn users_refresh_token_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.refresh_token.list")
    }
    pub fn users_refresh_token_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.refresh_token.delete")
    }
    pub fn users_risk_assessments_clear(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.risk_assessments.clear")
    }
    pub fn users_roles_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.roles.list")
    }
    pub fn users_roles_assign(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.roles.assign")
    }
    pub fn users_roles_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.roles.delete")
    }
    pub fn users_sessions_list(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.sessions.list")
    }
    pub fn users_sessions_delete(&self) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.sessions.delete")
    }
    pub fn users_effectivepermissions_sources_roles_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.effectivepermissions.sources.roles.list")
    }
    pub fn users_effectiveroles_sources_groups_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("users.effectiveroles.sources.groups.list")
    }
    pub fn verifiablecredentials_verification_templates_list(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.list")
    }
    pub fn verifiablecredentials_verification_templates_create(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.create")
    }
    pub fn verifiablecredentials_verification_templates_get(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.get")
    }
    pub fn verifiablecredentials_verification_templates_delete(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.delete")
    }
    pub fn verifiablecredentials_verification_templates_update(
        &self,
    ) -> Result<ManagementRequest<'_>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.update")
    }
}

impl<'a> RawManagementApi<'a> {
    pub fn actions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.list")
    }
    pub fn actions_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.create")
    }
    pub fn actions_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.get")
    }
    pub fn actions_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.delete")
    }
    pub fn actions_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.update")
    }
    pub fn actions_deploy(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.deploy")
    }
    pub fn actions_test(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.test")
    }
    pub fn branding_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.get")
    }
    pub fn branding_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.update")
    }
    pub fn client_grants_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("client_grants.list")
    }
    pub fn client_grants_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("client_grants.create")
    }
    pub fn client_grants_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("client_grants.get")
    }
    pub fn client_grants_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("client_grants.delete")
    }
    pub fn client_grants_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("client_grants.update")
    }
    pub fn clients_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.list")
    }
    pub fn clients_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.create")
    }
    pub fn clients_preview_cimd_metadata(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.preview_cimd_metadata")
    }
    pub fn clients_register_cimd_client(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.register_cimd_client")
    }
    pub fn clients_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.get")
    }
    pub fn clients_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.delete")
    }
    pub fn clients_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.update")
    }
    pub fn clients_rotate_secret(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.rotate_secret")
    }
    pub fn connection_profiles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.list")
    }
    pub fn connection_profiles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.create")
    }
    pub fn connection_profiles_list_templates(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.list_templates")
    }
    pub fn connection_profiles_get_template(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.get_template")
    }
    pub fn connection_profiles_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.get")
    }
    pub fn connection_profiles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.delete")
    }
    pub fn connection_profiles_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connection_profiles.update")
    }
    pub fn connections_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.list")
    }
    pub fn connections_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.create")
    }
    pub fn connections_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.get")
    }
    pub fn connections_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.delete")
    }
    pub fn connections_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.update")
    }
    pub fn connections_check_status(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.check_status")
    }
    pub fn custom_domains_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.list")
    }
    pub fn custom_domains_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.create")
    }
    pub fn custom_domains_get_default(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.get_default")
    }
    pub fn custom_domains_set_default(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.set_default")
    }
    pub fn custom_domains_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.get")
    }
    pub fn custom_domains_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.delete")
    }
    pub fn custom_domains_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.update")
    }
    pub fn custom_domains_test(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.test")
    }
    pub fn custom_domains_verify(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("custom_domains.verify")
    }
    pub fn device_credentials_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("device_credentials.list")
    }
    pub fn device_credentials_create_public_key(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("device_credentials.create_public_key")
    }
    pub fn device_credentials_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("device_credentials.delete")
    }
    pub fn email_templates_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("email_templates.create")
    }
    pub fn email_templates_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("email_templates.get")
    }
    pub fn email_templates_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("email_templates.set")
    }
    pub fn email_templates_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("email_templates.update")
    }
    pub fn event_streams_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("event_streams.list")
    }
    pub fn event_streams_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("event_streams.create")
    }
    pub fn event_streams_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("event_streams.get")
    }
    pub fn event_streams_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("event_streams.delete")
    }
    pub fn event_streams_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("event_streams.update")
    }
    pub fn event_streams_test(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("event_streams.test")
    }
    pub fn events_subscribe(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("events.subscribe")
    }
    pub fn flows_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.list")
    }
    pub fn flows_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.create")
    }
    pub fn flows_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.get")
    }
    pub fn flows_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.delete")
    }
    pub fn flows_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.update")
    }
    pub fn forms_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("forms.list")
    }
    pub fn forms_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("forms.create")
    }
    pub fn forms_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("forms.get")
    }
    pub fn forms_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("forms.delete")
    }
    pub fn forms_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("forms.update")
    }
    pub fn groups_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.list")
    }
    pub fn groups_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.get")
    }
    pub fn groups_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.delete")
    }
    pub fn hooks_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.list")
    }
    pub fn hooks_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.create")
    }
    pub fn hooks_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.get")
    }
    pub fn hooks_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.delete")
    }
    pub fn hooks_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.update")
    }
    pub fn jobs_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.get")
    }
    pub fn log_streams_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("log_streams.list")
    }
    pub fn log_streams_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("log_streams.create")
    }
    pub fn log_streams_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("log_streams.get")
    }
    pub fn log_streams_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("log_streams.delete")
    }
    pub fn log_streams_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("log_streams.update")
    }
    pub fn logs_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("logs.list")
    }
    pub fn logs_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("logs.get")
    }
    pub fn network_acls_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("network_acls.list")
    }
    pub fn network_acls_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("network_acls.create")
    }
    pub fn network_acls_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("network_acls.get")
    }
    pub fn network_acls_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("network_acls.set")
    }
    pub fn network_acls_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("network_acls.delete")
    }
    pub fn network_acls_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("network_acls.update")
    }
    pub fn organizations_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.list")
    }
    pub fn organizations_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.create")
    }
    pub fn organizations_get_by_name(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.get_by_name")
    }
    pub fn organizations_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.get")
    }
    pub fn organizations_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.delete")
    }
    pub fn organizations_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.update")
    }
    pub fn prompts_get_settings(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.get_settings")
    }
    pub fn prompts_update_settings(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.update_settings")
    }
    pub fn rate_limit_policies_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rate_limit_policies.list")
    }
    pub fn rate_limit_policies_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rate_limit_policies.create")
    }
    pub fn rate_limit_policies_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rate_limit_policies.get")
    }
    pub fn rate_limit_policies_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rate_limit_policies.delete")
    }
    pub fn rate_limit_policies_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rate_limit_policies.update")
    }
    pub fn refresh_tokens_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("refresh_tokens.list")
    }
    pub fn refresh_tokens_revoke(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("refresh_tokens.revoke")
    }
    pub fn refresh_tokens_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("refresh_tokens.get")
    }
    pub fn refresh_tokens_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("refresh_tokens.delete")
    }
    pub fn refresh_tokens_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("refresh_tokens.update")
    }
    pub fn resource_servers_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("resource_servers.list")
    }
    pub fn resource_servers_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("resource_servers.create")
    }
    pub fn resource_servers_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("resource_servers.get")
    }
    pub fn resource_servers_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("resource_servers.delete")
    }
    pub fn resource_servers_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("resource_servers.update")
    }
    pub fn roles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.list")
    }
    pub fn roles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.create")
    }
    pub fn roles_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.get")
    }
    pub fn roles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.delete")
    }
    pub fn roles_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.update")
    }
    pub fn rules_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules.list")
    }
    pub fn rules_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules.create")
    }
    pub fn rules_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules.get")
    }
    pub fn rules_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules.delete")
    }
    pub fn rules_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules.update")
    }
    pub fn rules_configs_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules_configs.list")
    }
    pub fn rules_configs_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules_configs.set")
    }
    pub fn rules_configs_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("rules_configs.delete")
    }
    pub fn self_service_profiles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("self_service_profiles.list")
    }
    pub fn self_service_profiles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("self_service_profiles.create")
    }
    pub fn self_service_profiles_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("self_service_profiles.get")
    }
    pub fn self_service_profiles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("self_service_profiles.delete")
    }
    pub fn self_service_profiles_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("self_service_profiles.update")
    }
    pub fn sessions_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("sessions.get")
    }
    pub fn sessions_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("sessions.delete")
    }
    pub fn sessions_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("sessions.update")
    }
    pub fn sessions_revoke(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("sessions.revoke")
    }
    pub fn stats_get_active_users_count(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("stats.get_active_users_count")
    }
    pub fn stats_get_daily(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("stats.get_daily")
    }
    pub fn supplemental_signals_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("supplemental_signals.get")
    }
    pub fn supplemental_signals_patch(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("supplemental_signals.patch")
    }
    pub fn tickets_verify_email(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("tickets.verify_email")
    }
    pub fn tickets_change_password(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("tickets.change_password")
    }
    pub fn token_exchange_profiles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("token_exchange_profiles.list")
    }
    pub fn token_exchange_profiles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("token_exchange_profiles.create")
    }
    pub fn token_exchange_profiles_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("token_exchange_profiles.get")
    }
    pub fn token_exchange_profiles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("token_exchange_profiles.delete")
    }
    pub fn token_exchange_profiles_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("token_exchange_profiles.update")
    }
    pub fn user_attribute_profiles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.list")
    }
    pub fn user_attribute_profiles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.create")
    }
    pub fn user_attribute_profiles_list_templates(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.list_templates")
    }
    pub fn user_attribute_profiles_get_template(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.get_template")
    }
    pub fn user_attribute_profiles_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.get")
    }
    pub fn user_attribute_profiles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.delete")
    }
    pub fn user_attribute_profiles_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_attribute_profiles.update")
    }
    pub fn user_blocks_list_by_identifier(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_blocks.list_by_identifier")
    }
    pub fn user_blocks_delete_by_identifier(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_blocks.delete_by_identifier")
    }
    pub fn user_blocks_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_blocks.list")
    }
    pub fn user_blocks_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_blocks.delete")
    }
    pub fn user_grants_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_grants.list")
    }
    pub fn user_grants_delete_by_user_id(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_grants.delete_by_user_id")
    }
    pub fn user_grants_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("user_grants.delete")
    }
    pub fn users_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.list")
    }
    pub fn users_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.create")
    }
    pub fn users_list_users_by_email(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.list_users_by_email")
    }
    pub fn users_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.get")
    }
    pub fn users_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.delete")
    }
    pub fn users_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.update")
    }
    pub fn users_regenerate_recovery_code(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.regenerate_recovery_code")
    }
    pub fn users_revoke_access(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.revoke_access")
    }
    pub fn actions_executions_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.executions.get")
    }
    pub fn actions_modules_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.list")
    }
    pub fn actions_modules_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.create")
    }
    pub fn actions_modules_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.get")
    }
    pub fn actions_modules_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.delete")
    }
    pub fn actions_modules_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.update")
    }
    pub fn actions_modules_list_actions(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.list_actions")
    }
    pub fn actions_modules_rollback(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.rollback")
    }
    pub fn actions_triggers_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.triggers.list")
    }
    pub fn actions_versions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.versions.list")
    }
    pub fn actions_versions_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.versions.get")
    }
    pub fn actions_versions_deploy(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.versions.deploy")
    }
    pub fn actions_modules_versions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.versions.list")
    }
    pub fn actions_modules_versions_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.versions.create")
    }
    pub fn actions_modules_versions_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.modules.versions.get")
    }
    pub fn actions_triggers_bindings_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.triggers.bindings.list")
    }
    pub fn actions_triggers_bindings_update_many(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("actions.triggers.bindings.update_many")
    }
    pub fn anomaly_blocks_check_ip(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("anomaly.blocks.check_ip")
    }
    pub fn anomaly_blocks_unblock_ip(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("anomaly.blocks.unblock_ip")
    }
    pub fn attackprotection_bot_detection_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.bot_detection.get")
    }
    pub fn attackprotection_bot_detection_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.bot_detection.update")
    }
    pub fn attackprotection_breached_password_detection_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.breached_password_detection.get")
    }
    pub fn attackprotection_breached_password_detection_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.breached_password_detection.update")
    }
    pub fn attackprotection_brute_force_protection_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.brute_force_protection.get")
    }
    pub fn attackprotection_brute_force_protection_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.brute_force_protection.update")
    }
    pub fn attackprotection_captcha_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.captcha.get")
    }
    pub fn attackprotection_captcha_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.captcha.update")
    }
    pub fn attackprotection_phone_provider_protection_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.phone_provider_protection.get")
    }
    pub fn attackprotection_phone_provider_protection_patch(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.phone_provider_protection.patch")
    }
    pub fn attackprotection_suspicious_ip_throttling_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.suspicious_ip_throttling.get")
    }
    pub fn attackprotection_suspicious_ip_throttling_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("attackprotection.suspicious_ip_throttling.update")
    }
    pub fn branding_templates_get_universal_login(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.templates.get_universal_login")
    }
    pub fn branding_templates_update_universal_login(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.templates.update_universal_login")
    }
    pub fn branding_templates_delete_universal_login(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.templates.delete_universal_login")
    }
    pub fn branding_themes_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.themes.create")
    }
    pub fn branding_themes_get_default(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.themes.get_default")
    }
    pub fn branding_themes_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.themes.get")
    }
    pub fn branding_themes_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.themes.delete")
    }
    pub fn branding_themes_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.themes.update")
    }
    pub fn branding_phone_providers_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.providers.list")
    }
    pub fn branding_phone_providers_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.providers.create")
    }
    pub fn branding_phone_providers_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.providers.get")
    }
    pub fn branding_phone_providers_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.providers.delete")
    }
    pub fn branding_phone_providers_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.providers.update")
    }
    pub fn branding_phone_providers_test(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.providers.test")
    }
    pub fn branding_phone_templates_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.list")
    }
    pub fn branding_phone_templates_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.create")
    }
    pub fn branding_phone_templates_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.get")
    }
    pub fn branding_phone_templates_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.delete")
    }
    pub fn branding_phone_templates_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.update")
    }
    pub fn branding_phone_templates_reset(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.reset")
    }
    pub fn branding_phone_templates_test(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("branding.phone.templates.test")
    }
    pub fn clientgrants_organizations_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clientgrants.organizations.list")
    }
    pub fn clients_connections_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.connections.get")
    }
    pub fn clients_credentials_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.credentials.list")
    }
    pub fn clients_credentials_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.credentials.create")
    }
    pub fn clients_credentials_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.credentials.get")
    }
    pub fn clients_credentials_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.credentials.delete")
    }
    pub fn clients_credentials_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("clients.credentials.update")
    }
    pub fn connections_clients_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.clients.get")
    }
    pub fn connections_clients_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.clients.update")
    }
    pub fn connections_directory_provisioning_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.list")
    }
    pub fn connections_directory_provisioning_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.get")
    }
    pub fn connections_directory_provisioning_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.create")
    }
    pub fn connections_directory_provisioning_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.delete")
    }
    pub fn connections_directory_provisioning_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.update")
    }
    pub fn connections_directory_provisioning_get_default_mapping(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.get_default_mapping")
    }
    pub fn connections_directory_provisioning_list_synchronized_groups(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.list_synchronized_groups")
    }
    pub fn connections_directory_provisioning_set(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directory_provisioning.set")
    }
    pub fn connections_keys_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.keys.get")
    }
    pub fn connections_keys_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.keys.create")
    }
    pub fn connections_keys_rotate(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.keys.rotate")
    }
    pub fn connections_scim_configuration_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scim_configuration.list")
    }
    pub fn connections_scim_configuration_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scim_configuration.get")
    }
    pub fn connections_scim_configuration_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scim_configuration.create")
    }
    pub fn connections_scim_configuration_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scim_configuration.delete")
    }
    pub fn connections_scim_configuration_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scim_configuration.update")
    }
    pub fn connections_scim_configuration_get_default_mapping(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scim_configuration.get_default_mapping")
    }
    pub fn connections_users_delete_by_email(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.users.delete_by_email")
    }
    pub fn connections_directoryprovisioning_synchronizations_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.directoryprovisioning.synchronizations.create")
    }
    pub fn connections_scimconfiguration_tokens_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scimconfiguration.tokens.get")
    }
    pub fn connections_scimconfiguration_tokens_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scimconfiguration.tokens.create")
    }
    pub fn connections_scimconfiguration_tokens_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("connections.scimconfiguration.tokens.delete")
    }
    pub fn emails_provider_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("emails.provider.get")
    }
    pub fn emails_provider_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("emails.provider.create")
    }
    pub fn emails_provider_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("emails.provider.delete")
    }
    pub fn emails_provider_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("emails.provider.update")
    }
    pub fn eventstreams_deliveries_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("eventstreams.deliveries.list")
    }
    pub fn eventstreams_deliveries_get_history(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("eventstreams.deliveries.get_history")
    }
    pub fn eventstreams_redeliveries_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("eventstreams.redeliveries.create")
    }
    pub fn eventstreams_redeliveries_create_by_id(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("eventstreams.redeliveries.create_by_id")
    }
    pub fn flows_executions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.executions.list")
    }
    pub fn flows_executions_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.executions.get")
    }
    pub fn flows_executions_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.executions.delete")
    }
    pub fn flows_vault_connections_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.vault.connections.list")
    }
    pub fn flows_vault_connections_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.vault.connections.create")
    }
    pub fn flows_vault_connections_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.vault.connections.get")
    }
    pub fn flows_vault_connections_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.vault.connections.delete")
    }
    pub fn flows_vault_connections_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("flows.vault.connections.update")
    }
    pub fn groups_members_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.members.get")
    }
    pub fn groups_roles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.roles.list")
    }
    pub fn groups_roles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.roles.create")
    }
    pub fn groups_roles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("groups.roles.delete")
    }
    pub fn guardian_enrollments_create_ticket(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.enrollments.create_ticket")
    }
    pub fn guardian_enrollments_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.enrollments.get")
    }
    pub fn guardian_enrollments_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.enrollments.delete")
    }
    pub fn guardian_factors_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.list")
    }
    pub fn guardian_factors_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.set")
    }
    pub fn guardian_policies_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.policies.list")
    }
    pub fn guardian_policies_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.policies.set")
    }
    pub fn guardian_factors_phone_get_message_types(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.get_message_types")
    }
    pub fn guardian_factors_phone_set_message_types(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.set_message_types")
    }
    pub fn guardian_factors_phone_get_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.get_twilio_provider")
    }
    pub fn guardian_factors_phone_set_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.set_twilio_provider")
    }
    pub fn guardian_factors_phone_get_selected_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.get_selected_provider")
    }
    pub fn guardian_factors_phone_set_provider(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.set_provider")
    }
    pub fn guardian_factors_phone_get_templates(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.get_templates")
    }
    pub fn guardian_factors_phone_set_templates(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.phone.set_templates")
    }
    pub fn guardian_factors_push_notification_get_apns_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.get_apns_provider")
    }
    pub fn guardian_factors_push_notification_set_apns_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_apns_provider")
    }
    pub fn guardian_factors_push_notification_update_apns_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_apns_provider")
    }
    pub fn guardian_factors_push_notification_set_fcm_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_fcm_provider")
    }
    pub fn guardian_factors_push_notification_update_fcm_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_fcm_provider")
    }
    pub fn guardian_factors_push_notification_set_fcmv1_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_fcmv1_provider")
    }
    pub fn guardian_factors_push_notification_update_fcmv1_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_fcmv1_provider")
    }
    pub fn guardian_factors_push_notification_get_sns_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.get_sns_provider")
    }
    pub fn guardian_factors_push_notification_set_sns_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_sns_provider")
    }
    pub fn guardian_factors_push_notification_update_sns_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.update_sns_provider")
    }
    pub fn guardian_factors_push_notification_get_selected_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.get_selected_provider")
    }
    pub fn guardian_factors_push_notification_set_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.push_notification.set_provider")
    }
    pub fn guardian_factors_sms_get_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.sms.get_twilio_provider")
    }
    pub fn guardian_factors_sms_set_twilio_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.sms.set_twilio_provider")
    }
    pub fn guardian_factors_sms_get_selected_provider(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.sms.get_selected_provider")
    }
    pub fn guardian_factors_sms_set_provider(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.sms.set_provider")
    }
    pub fn guardian_factors_sms_get_templates(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.sms.get_templates")
    }
    pub fn guardian_factors_sms_set_templates(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.sms.set_templates")
    }
    pub fn guardian_factors_duo_settings_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.duo.settings.get")
    }
    pub fn guardian_factors_duo_settings_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.duo.settings.set")
    }
    pub fn guardian_factors_duo_settings_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("guardian.factors.duo.settings.update")
    }
    pub fn hooks_secrets_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.secrets.get")
    }
    pub fn hooks_secrets_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.secrets.create")
    }
    pub fn hooks_secrets_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.secrets.delete")
    }
    pub fn hooks_secrets_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("hooks.secrets.update")
    }
    pub fn jobs_errors_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.errors.get")
    }
    pub fn jobs_users_exports_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.users_exports.create")
    }
    pub fn jobs_users_imports_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.users_imports.create")
    }
    pub fn jobs_users_imports_create_post(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.users_imports.create_post")
    }
    pub fn jobs_users_imports_create_post_3(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.users_imports.create_post_3")
    }
    pub fn jobs_users_imports_create_post_4(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.users_imports.create_post_4")
    }
    pub fn jobs_users_imports_create_post_5(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.users_imports.create_post_5")
    }
    pub fn jobs_verification_email_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("jobs.verification_email.create")
    }
    pub fn keys_custom_signing_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.custom_signing.get")
    }
    pub fn keys_custom_signing_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.custom_signing.set")
    }
    pub fn keys_custom_signing_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.custom_signing.delete")
    }
    pub fn keys_encryption_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.list")
    }
    pub fn keys_encryption_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.create")
    }
    pub fn keys_encryption_rekey(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.rekey")
    }
    pub fn keys_encryption_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.get")
    }
    pub fn keys_encryption_import_(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.import_")
    }
    pub fn keys_encryption_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.delete")
    }
    pub fn keys_encryption_create_public_wrapping_key(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.encryption.create_public_wrapping_key")
    }
    pub fn keys_signing_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.signing.list")
    }
    pub fn keys_signing_rotate(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.signing.rotate")
    }
    pub fn keys_signing_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.signing.get")
    }
    pub fn keys_signing_revoke(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("keys.signing.revoke")
    }
    pub fn organizations_client_grants_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.client_grants.list")
    }
    pub fn organizations_client_grants_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.client_grants.create")
    }
    pub fn organizations_client_grants_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.client_grants.delete")
    }
    pub fn organizations_connections_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.connections.list")
    }
    pub fn organizations_connections_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.connections.create")
    }
    pub fn organizations_connections_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.connections.get")
    }
    pub fn organizations_connections_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.connections.delete")
    }
    pub fn organizations_connections_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.connections.update")
    }
    pub fn organizations_discovery_domains_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.discovery_domains.list")
    }
    pub fn organizations_discovery_domains_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.discovery_domains.create")
    }
    pub fn organizations_discovery_domains_get_by_name(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.discovery_domains.get_by_name")
    }
    pub fn organizations_discovery_domains_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.discovery_domains.get")
    }
    pub fn organizations_discovery_domains_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.discovery_domains.delete")
    }
    pub fn organizations_discovery_domains_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.discovery_domains.update")
    }
    pub fn organizations_enabled_connections_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.enabled_connections.list")
    }
    pub fn organizations_enabled_connections_add(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.enabled_connections.add")
    }
    pub fn organizations_enabled_connections_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.enabled_connections.get")
    }
    pub fn organizations_enabled_connections_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.enabled_connections.delete")
    }
    pub fn organizations_enabled_connections_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.enabled_connections.update")
    }
    pub fn organizations_groups_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.groups.list")
    }
    pub fn organizations_invitations_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.invitations.list")
    }
    pub fn organizations_invitations_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.invitations.create")
    }
    pub fn organizations_invitations_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.invitations.get")
    }
    pub fn organizations_invitations_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.invitations.delete")
    }
    pub fn organizations_members_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.list")
    }
    pub fn organizations_members_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.create")
    }
    pub fn organizations_members_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.delete")
    }
    pub fn organizations_groups_roles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.groups.roles.list")
    }
    pub fn organizations_groups_roles_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.groups.roles.create")
    }
    pub fn organizations_groups_roles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.groups.roles.delete")
    }
    pub fn organizations_members_effective_roles_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.effective_roles.list")
    }
    pub fn organizations_members_roles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.roles.list")
    }
    pub fn organizations_members_roles_assign(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.roles.assign")
    }
    pub fn organizations_members_roles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.roles.delete")
    }
    pub fn organizations_members_effectiveroles_sources_groups_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("organizations.members.effectiveroles.sources.groups.list")
    }
    pub fn prompts_custom_text_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.custom_text.get")
    }
    pub fn prompts_custom_text_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.custom_text.set")
    }
    pub fn prompts_partials_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.partials.get")
    }
    pub fn prompts_partials_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.partials.set")
    }
    pub fn prompts_rendering_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.rendering.list")
    }
    pub fn prompts_rendering_bulk_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.rendering.bulk_update")
    }
    pub fn prompts_rendering_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.rendering.get")
    }
    pub fn prompts_rendering_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("prompts.rendering.update")
    }
    pub fn riskassessments_settings_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("riskassessments.settings.get")
    }
    pub fn riskassessments_settings_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("riskassessments.settings.update")
    }
    pub fn riskassessments_settings_new_device_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("riskassessments.settings.new_device.get")
    }
    pub fn riskassessments_settings_new_device_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("riskassessments.settings.new_device.update")
    }
    pub fn roles_groups_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.groups.get")
    }
    pub fn roles_groups_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.groups.create")
    }
    pub fn roles_groups_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.groups.delete")
    }
    pub fn roles_permissions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.permissions.list")
    }
    pub fn roles_permissions_add(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.permissions.add")
    }
    pub fn roles_permissions_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.permissions.delete")
    }
    pub fn roles_users_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.users.list")
    }
    pub fn roles_users_assign(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("roles.users.assign")
    }
    pub fn selfserviceprofiles_custom_text_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("selfserviceprofiles.custom_text.list")
    }
    pub fn selfserviceprofiles_custom_text_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("selfserviceprofiles.custom_text.set")
    }
    pub fn selfserviceprofiles_sso_ticket_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("selfserviceprofiles.sso_ticket.create")
    }
    pub fn selfserviceprofiles_sso_ticket_revoke(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("selfserviceprofiles.sso_ticket.revoke")
    }
    pub fn tenants_settings_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("tenants.settings.get")
    }
    pub fn tenants_settings_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("tenants.settings.update")
    }
    pub fn users_authentication_methods_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.list")
    }
    pub fn users_authentication_methods_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.create")
    }
    pub fn users_authentication_methods_set(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.set")
    }
    pub fn users_authentication_methods_delete_all(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.delete_all")
    }
    pub fn users_authentication_methods_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.get")
    }
    pub fn users_authentication_methods_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.delete")
    }
    pub fn users_authentication_methods_update(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authentication_methods.update")
    }
    pub fn users_authenticators_delete_all(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.authenticators.delete_all")
    }
    pub fn users_connected_accounts_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.connected_accounts.list")
    }
    pub fn users_effective_permissions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.effective_permissions.list")
    }
    pub fn users_effective_roles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.effective_roles.list")
    }
    pub fn users_enrollments_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.enrollments.get")
    }
    pub fn users_federated_connections_tokensets_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.federated_connections_tokensets.list")
    }
    pub fn users_federated_connections_tokensets_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.federated_connections_tokensets.delete")
    }
    pub fn users_groups_get(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.groups.get")
    }
    pub fn users_identities_link(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.identities.link")
    }
    pub fn users_identities_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.identities.delete")
    }
    pub fn users_logs_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.logs.list")
    }
    pub fn users_multifactor_invalidate_remember_browser(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.multifactor.invalidate_remember_browser")
    }
    pub fn users_multifactor_delete_provider(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.multifactor.delete_provider")
    }
    pub fn users_organizations_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.organizations.list")
    }
    pub fn users_permissions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.permissions.list")
    }
    pub fn users_permissions_create(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.permissions.create")
    }
    pub fn users_permissions_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.permissions.delete")
    }
    pub fn users_refresh_token_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.refresh_token.list")
    }
    pub fn users_refresh_token_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.refresh_token.delete")
    }
    pub fn users_risk_assessments_clear(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.risk_assessments.clear")
    }
    pub fn users_roles_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.roles.list")
    }
    pub fn users_roles_assign(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.roles.assign")
    }
    pub fn users_roles_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.roles.delete")
    }
    pub fn users_sessions_list(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.sessions.list")
    }
    pub fn users_sessions_delete(&self) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.sessions.delete")
    }
    pub fn users_effectivepermissions_sources_roles_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.effectivepermissions.sources.roles.list")
    }
    pub fn users_effectiveroles_sources_groups_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("users.effectiveroles.sources.groups.list")
    }
    pub fn verifiablecredentials_verification_templates_list(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.list")
    }
    pub fn verifiablecredentials_verification_templates_create(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.create")
    }
    pub fn verifiablecredentials_verification_templates_get(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.get")
    }
    pub fn verifiablecredentials_verification_templates_delete(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.delete")
    }
    pub fn verifiablecredentials_verification_templates_update(
        &self,
    ) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.request("verifiablecredentials.verification.templates.update")
    }
}
