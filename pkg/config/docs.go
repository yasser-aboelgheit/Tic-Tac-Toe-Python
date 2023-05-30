/* Package config-lib helps you read all the configurations by very simple approach of give it empty struct
with predefined mapstructures in the struct tags, and default function to call if there is any
default value need to be settled.

config-lib comes handy from infra prespictive because we can change how we deploy the apps env secrets
without affecting the development of any service, also it applies for all services.

* It's like a wrapper to read the env variables in our own way that follows
* 12 factor.
* as start if you provided filePath the envs in the in the file should not contain the prefix
* unlike the environment variable you must use the prefix with `_` underscore, for example if the env `key1` and prefix is `service1` we will search for `key1_service1`
 */
package config
