# \DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Cash**](DefaultApi.md#Cash) | **Post** /cash | 
[**Dig**](DefaultApi.md#Dig) | **Post** /dig | 
[**ExploreArea**](DefaultApi.md#ExploreArea) | **Post** /explore | 
[**GetBalance**](DefaultApi.md#GetBalance) | **Get** /balance | 
[**HealthCheck**](DefaultApi.md#HealthCheck) | **Get** /health-check | 
[**IssueLicense**](DefaultApi.md#IssueLicense) | **Post** /licenses | 
[**ListLicenses**](DefaultApi.md#ListLicenses) | **Get** /licenses | 


# **Cash**
> Wallet Cash(ctx, args)


Exchange provided treasure for money.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **args** | [**Treasure**](Treasure.md)| Treasure for exchange. | 

### Return type

[**Wallet**](wallet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Dig**
> TreasureList Dig(ctx, args)


Dig at given point and depth, returns found treasures.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **args** | [**Dig**](Dig.md)| License, place and depth to dig. | 

### Return type

[**TreasureList**](treasureList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ExploreArea**
> Report ExploreArea(ctx, args)


Returns amount of treasures in the provided area at full depth.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **args** | [**Area**](Area.md)| Area to be explored. | 

### Return type

[**Report**](report.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetBalance**
> Balance GetBalance(ctx, )


Returns a current balance.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**Balance**](balance.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **HealthCheck**
> map[string]interface{} HealthCheck(ctx, )


Returns 200 if service works okay.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**map[string]interface{}**](interface{}.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **IssueLicense**
> License IssueLicense(ctx, optional)


Issue a new license.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DefaultApiIssueLicenseOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiIssueLicenseOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **args** | [**optional.Interface of Wallet**](Wallet.md)| Amount of money to spend for a license. Empty array for get free license. Maximum 10 active licenses | 

### Return type

[**License**](license.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListLicenses**
> LicenseList ListLicenses(ctx, )


Returns a list of issued licenses.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**LicenseList**](licenseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

