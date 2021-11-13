/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
const fs = require('fs'); // 파일 읽기 쓰기 js 기본패키지
const path = require('path'); // 경로명에 해당하는 문자열처리 기본패키지
// _dirname: nodejsprogram (enrollAdmin.js 위치)
const ccpPath = path.resolve(__dirname, 'connection.json'); // ~dev/sky-project/connection.json
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON); // Json Unmarshal

async function main() {
    try {

        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.example.com'].url;
        const ca = new FabricCAServices(caURL);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet'); // ~dev/sky-project/Wallet
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (adminExists) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' }); // ca.example.com 컨테이너 -> sh: fabric-ca-server admin:adminpw
        const identity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
        wallet.import('admin', identity);
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        process.exit(1);
    }
}

main();
