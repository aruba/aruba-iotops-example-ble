

function decode(address, addressType, advType, elements)
    local handle = BLEData:new()

    MANUFACTURER_DATA_TYPE = 0xff

    APPLE_COMPANY_ID = 0x004c
    PROXIMITY_BCN_TYPE_VAL0 = 0x02
    PROXIMITY_BCN_TYPE_VAL1 = 0x15


    -- First, search manufacturer data
    local manufacturer_data = elements[MANUFACTURER_DATA_TYPE]
    if manufacturer_data ~= nil then

        -- check length
        if #manufacturer_data < 25 then
            return handle
        end  

        -- get the manufacturer id
        local company_id = string.unpack("<I2", (string.sub(manufacturer_data, 1, 2)))

        if company_id == APPLE_COMPANY_ID then
            local beacon_type0 = string.unpack("<I1", string.sub(manufacturer_data, 3, 3))
            if beacon_type0 ~= PROXIMITY_BCN_TYPE_VAL0 then
                return handle
            end

            local beacon_type1 = string.unpack("<I1", string.sub(manufacturer_data, 4, 4))
            if beacon_type1 ~= PROXIMITY_BCN_TYPE_VAL1 then
                return handle
            end

            local proximity_uuid = nil
            local proximity_major = nil
            local proximity_minor = nil
            local proximity_measured_power = nil
            local extra = nil

            proximity_uuid = string.sub(manufacturer_data, 5, 20)
            proximity_major = string.unpack(">I2", string.sub(manufacturer_data, 21, 22))
            proximity_minor = string.unpack(">I2", string.sub(manufacturer_data, 23, 24))
            proximity_measured_power = string.unpack(">i1", string.sub(manufacturer_data, 25, 25))
            if #manufacturer_data == 26 then
                extra = string.unpack(">I1", string.sub(manufacturer_data, 26, 26))
            end

            -- all good. populate return handle
            handle:setDeviceClass("iBeacon")
            handle:setIBeaconUUID(proximity_uuid)
            handle:setIBeaconMajor(proximity_major)
            handle:setIBeaconMinor(proximity_minor)
            handle:setIBeaconPower(proximity_measured_power)
            if extra ~= nil then
                handle:setIBeaconExtra(extra)
            end
        end
    end

    return handle
end
