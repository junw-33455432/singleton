/*
 * Copyright 2019-2020 VMware, Inc.
 * SPDX-License-Identifier: EPL-2.0
 */
package com.vmware.vip.i18n.api.v1.authentication;

import javax.servlet.http.HttpServletRequest;

import org.apache.commons.lang3.StringUtils;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import com.vmware.vip.api.rest.API;
import com.vmware.vip.api.rest.APIOperation;
import com.vmware.vip.api.rest.APIParamName;
import com.vmware.vip.api.rest.APIParamValue;
import com.vmware.vip.api.rest.APIV1;
import com.vmware.vip.common.i18n.dto.AuthenKeyDTO;
import com.vmware.vip.common.i18n.dto.response.AthenticationResponseDTO;
import com.vmware.vip.common.i18n.status.APIResponseStatus;
import com.vmware.vip.i18n.api.v1.utils.KeyService;

import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;

/**
 * Provide RESTful API to manipulate the key for product.
 *
 */
@RestController
public class SecurityKeyAPI {


    /**
     * Generate the key by productName,version,and userID.
     * and return it to product.
     * <p>
     * If the request parameters have error, return Bad Request to product.
     * <p>
     * else if the key is generated successfully, return the key to product.
     * <p>
     * else return the Internal Server Error to product.
     *
     * @param productName
     *        The name of product.
     * @param version
     *        The release version of product.
     * @param userID
     *        User identifier randomly generated by product
     * @param request
     *        Extends the ServletRequest interface to provide request information for HTTP servlets.
     * @return APIResponseDTO 
     *         The object which represents response status.
     */
    /*
     * Temporarily disabling this unused API by commenting out the following annotations. - 09May19 
    */
    public AthenticationResponseDTO getKey(
            @ApiParam(name = APIParamName.PRODUCT_NAME, required = true, value = APIParamValue.PRODUCT_NAME) @RequestParam(value = APIParamName.PRODUCT_NAME, required = true) String productName,
            @ApiParam(name = APIParamName.VERSION, required = true, value = APIParamName.VERSION) @RequestParam(value = APIParamName.VERSION, required = true) String version,
            @ApiParam(name = APIParamName.USER_ID, required = true, value = APIParamValue.USERID) @RequestParam(value = APIParamName.USER_ID, required = true) String userID,
            HttpServletRequest request) {
        AthenticationResponseDTO athenticationResponseDTO = new AthenticationResponseDTO();
        if (StringUtils.isNotEmpty(productName) && StringUtils.isNotEmpty(version)
                && StringUtils.isNotEmpty(userID)) {
            AuthenKeyDTO keyDTO = new AuthenKeyDTO();
            keyDTO.setProductName(productName);
            keyDTO.setVersion(version);
            keyDTO.setUserID(userID);
            keyDTO = KeyService.generateKey(keyDTO);
            if (null == keyDTO) {
                athenticationResponseDTO.setResponse(APIResponseStatus.INTERNAL_SERVER_ERROR);
                athenticationResponseDTO.setData("");
            } else {
                athenticationResponseDTO.setResponse(APIResponseStatus.OK);
                athenticationResponseDTO.setData(keyDTO);
            }
        } else {
            athenticationResponseDTO.setResponse(APIResponseStatus.BAD_REQUEST);
            athenticationResponseDTO.setData("");
        }
        return athenticationResponseDTO;
    }
}